package cron

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nvg14/hark/internal/database"
)

func RunUploadS3(hour, minute,seconds, intervalInSeconds int,bucket string,uploader *s3manager.Uploader,folders []string ) {
	now := time.Now()
	fmt.Println("Now:",now)
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0,0, 0, time.Local)
	fmt.Println("next-tick before if:",nextTick)
	if nextTick.Before(now) {
		nextTick = nextTick.Add(time.Duration(intervalInSeconds) * time.Second)
	}
	fmt.Println("next-tick:",nextTick)
	nextcron := time.Until(nextTick)
	fmt.Println("next-cron:",nextcron)
	t := time.NewTimer(nextcron)

	defer t.Stop()
	for {
		select {
		case <-t.C:
			fmt.Println("hi")
			for _,folder :=range folders{
				now := time.Now()
				before := now.Add(-1 * time.Hour)
				filepath := before.Format("2006-01-02_15")+ ".txt"
				err := database.UploadS3(bucket,uploader,folder,filepath)
				if err != nil{
					fmt.Println(err)
					fmt.Println("noooooooooo")
				}
				os.Remove(database.DirectoryPath+folder+"/"+filepath)
			}
			t.Reset(time.Duration(intervalInSeconds) * time.Second)
		}
	}
}