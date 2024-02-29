package poststore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type PostStore struct {
	cli *api.Client
}

func New() (*PostStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)

	if err != nil {
		return nil, err
	}

	return &PostStore{
		cli: client,
	}, nil
}

func (ps *PostStore) Get(id string, version string) ([]*Config, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	posts := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (ps *PostStore) GetAll() ([]*Config, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *PostStore) Delete(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *PostStore) Post(config *Config) (*Config, error) {
	kv := ps.cli.KV()

	sid, rid := generateKey(config.Version, config.Labels)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (ps *PostStore) GetPostsByLabels(id string, version string, labels string) ([]*RequestPost, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, labels), nil)
	if err != nil {
		return nil, err
	}

	posts := []*RequestPost{}

	for _, pair := range data {
		post := &RequestPost{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (ps *PostStore) PostGroup(group *Group) (*Group, error) {
	kv := ps.cli.KV()

	sid, rid := generateKey(group.Version, group.Labels)
	group.Id = rid

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}
func (ps *PostStore) GetAllGroups() ([]*Group, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil) //lista grupa
	if err != nil {
		return nil, err
	}

	posts := []*Group{}
	for _, pair := range data {
		post := &Group{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *PostStore) GetGroupById(id string, version string) ([]*Group, error) {
	kv := ps.cli.KV()

	listGrupa, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	groups := []*Group{}
	for _, pair := range listGrupa {
		post := &Group{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		if post.Id == id {
			if post.Version == version {
				groups = append(groups, post)

			}
		}

	}
	return groups, nil
}

func (ps *PostStore) DeleteGroup(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *PostStore) AddConfigToGroup(groupId string, configId string) error {
	kv := ps.cli.KV()

	configs, err := ps.Get(configId, "")
	if err != nil {
		return err
	}
	config := *configs[0]

	groups, err := ps.GetGroupById(groupId, "")
	if err != nil {
		return err
	}
	group := *groups[0]

	for _, c := range group.Configs {
		if c.Id == configId {
			return fmt.Errorf("Config with ID %s already exists in group with ID %s and version %s", configId, groupId, config.Version)
		}
	}

	group.Configs = append(group.Configs, config)

	data, err := json.Marshal(group)
	if err != nil {
		return err
	}

	sid, _ := generateKey(group.Id, group.Version)
	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}
