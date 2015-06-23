package model

import (
	"element"
	"errors"
	"sync"
)

type myMusicModel struct {
}

var myModel *myMusicModel = nil
var myMusicModelOnce sync.Once

func MyMusicModelInstance() *myMusicModel {
	myMusicModelOnce.Do(func() {
		myModel = &myMusicModel{}
	})

	return myModel
}

func (this *myMusicModel) InsertMusic(musicInfo *element.MusicInfo) (int, error) {
	stmt, err := DatabaseInstance().DB.Prepare("insert into localmusic(musicid, path) VALUES(?, ?)")
	_, err = stmt.Exec(musicInfo.MusicId, musicInfo.MusicUUID)
	return musicInfo.MusicId, err
}

func (this *myMusicModel) FetchMusicInfo(musicInfo *element.MusicInfo) error {
	rows, err := DatabaseInstance().DB.Query("select path from baidumusic where musicid = ?", musicInfo.MusicId)
	if err != nil {
		return err
	}
	if rows.Next() {
		err := rows.Scan(&musicInfo.MusicUUID)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("No Data")
	}
}

func (this *myMusicModel) DeleteMusic(musicId int) error {
	stmt, err := DatabaseInstance().DB.Prepare("delete from localmusic where musicid = ?")
	_, err = stmt.Exec(musicId)
	return err
}