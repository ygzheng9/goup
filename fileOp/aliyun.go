package main

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func getFileList(product string) ([]string, error) {
	result := []string{}
	// 建立连接
	client, err := oss.New("http://oss-cn-hangzhou.aliyuncs.com", "LTAIJuSjXPx3B35m", "5nz7Blk9nIr2R11Tss7FOVU5pk5cPJ")
	if err != nil {
		return result, err
	}

	// 设置 bucket
	bucket, err := client.Bucket("mfrwxoss")
	if err != nil {
		return result, err
	}

	// 根据 product 找到对应的文件名
	prefix := "300Wx300H/" + product
	lsRes, err := bucket.ListObjects(oss.Prefix(prefix))
	if err != nil {
		return result, err
	}

	for _, object := range lsRes.Objects {
		result = append(result, object.Key)
	}

	return result, nil
}

func connectAli() {
	client, err := oss.New("http://oss-cn-hangzhou.aliyuncs.com", "LTAIJuSjXPx3B35m", "5nz7Blk9nIr2R11Tss7FOVU5pk5cPJ")
	if err != nil {
		handleError(err)
	}

	// lsRes, err := client.ListBuckets()
	// if err != nil {
	// 	handleError(err)
	// }
	// for _, bucket := range lsRes.Buckets {
	// 	fmt.Println("bucket:", bucket.Name)
	// }

	bucket, err := client.Bucket("mfrwxoss")
	if err != nil {
		handleError(err)
	}

	lsRes, err := bucket.ListObjects(oss.Prefix("300Wx300H/YVA03-01-0004"))
	if err != nil {
		handleError(err)
	}
	for _, object := range lsRes.Objects {
		fmt.Println("Object:", object.Key)
	}

	// 通过循环， 获取所有的文件名
	// marker := oss.Marker("")
	// for {
	// 	lsRes, err := bucket.ListObjects(oss.Prefix("300Wx300H"), marker)
	// 	if err != nil {
	// 		handleError(err)
	// 	}
	// 	marker = oss.Marker(lsRes.NextMarker)

	// 	// fmt.Println("Objects:", lsRes.Objects)

	// 	for _, object := range lsRes.Objects {
	// 		fmt.Println("Object:", object.Key)
	// 	}

	// 	if !lsRes.IsTruncated {
	// 		break
	// 	}
	// }
}

func setCors() {
	client, err := oss.New("http://oss-cn-hangzhou.aliyuncs.com", "LTAIJuSjXPx3B35m", "5nz7Blk9nIr2R11Tss7FOVU5pk5cPJ")
	if err != nil {
		handleError(err)
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		handleError(err)
	}
	for _, bucket := range lsRes.Buckets {
		fmt.Println("bucket:", bucket.Name)
	}

	// 只有 bucket owner 才有权限修改
	// rule1 := oss.CORSRule{
	// 	AllowedOrigin: []string{"*"},
	// 	AllowedMethod: []string{"PUT", "GET"},
	// 	AllowedHeader: []string{},
	// 	ExposeHeader:  []string{},
	// 	MaxAgeSeconds: 200,
	// }
	// 查看Bucket ACL
	// aclRes, err := client.GetBucketACL("mfrwxoss")
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println("Bucket ACL:", aclRes.ACL)

	// err = client.SetBucketCORS("mfrwxoss", []oss.CORSRule{rule1})
	// if err != nil {
	// 	handleError(err)
	// }

	corsRes, err := client.GetBucketCORS("mfrwxoss")
	if err != nil {
		handleError(err)
	}
	fmt.Println("Bucket CORS:", corsRes.CORSRules)
}
