package internal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type CloudianAPI struct {
	client *resty.Client
}

func NewCloudianAPI(url string, username string, password string) CloudianAPI {
	return CloudianAPI{
		client: resty.New().
			SetBaseURL(url).
			SetBasicAuth(username, password).
			SetLogger(logrus.StandardLogger()),
	}
}

type UserBuckets struct {
	UserID string
	Bucket string
}

func (c CloudianAPI) GetBuckets(groupId string) ([]UserBuckets, error) {
	var bucketResponse []struct {
		UserID  string `json:"userId"`
		Buckets []struct {
			BucketName string `json:"bucketName"`
		} `json:"buckets"`
	}
	if r, err := c.client.R().
		SetResult(&bucketResponse).
		SetQueryParam("groupId", groupId).
		Get("/system/bucketlist"); err != nil {
		return []UserBuckets{}, err
	} else {
		if r.IsError() {
			return []UserBuckets{}, fmt.Errorf("error fetching buckets for group %s: (%s) %s", groupId, r.Status(), r.Body())
		}
		var result []UserBuckets
		for _, g := range bucketResponse {
			for _, bucket := range g.Buckets {
				result = append(result, UserBuckets{
					Bucket: bucket.BucketName,
					UserID: g.UserID,
				})
			}
		}
		return result, nil
	}
}

func (c CloudianAPI) GetGroups() ([]string, error) {
	var groupsResponse []struct {
		GroupID string `json:"groupId"`
	}
	if r, err := c.client.R().
		SetResult(&groupsResponse).
		Get("/group/list"); err != nil {
		return []string{}, err
	} else {
		if r.IsError() {
			return []string{}, fmt.Errorf("error fetching groups: (%s) %s", r.Status(), r.Body())
		}
		var result []string
		for _, g := range groupsResponse {
			result = append(result, g.GroupID)
		}
		return result, nil
	}
}

func (c CloudianAPI) GetBucketSize(groupID string, userID string, bucket string) (int64, error) {
	var bucketByteCount int64
	if r, err := c.client.R().
		SetQueryParam("groupId", groupID).
		SetQueryParam("userId", userID).
		SetQueryParam("bucketName", bucket).
		SetResult(&bucketByteCount).
		Get("/system/bytecount"); err != nil {
		return 0, err
	} else {
		if r.IsError() {
			return 0, fmt.Errorf("error fetching bucket size from Cloudian group %s, user %s, bucket %s: (%s) %s", groupID, userID, bucket, r.Status(), r.Body())
		}
		return bucketByteCount, nil
	}
}
