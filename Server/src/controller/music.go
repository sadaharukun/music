package controller

import (
	"bufio"
	"element"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"model"
	"music"
	"music/player"
	"net/http"
	"strconv"
	"token"
)

type MusicElementPacket struct {
	MusicID     int    `json:"id"`
	MusicName   string `json:"musicname"`
	ArtistName  string `json:"artist"`
	AlbumName   string `json:"albumname"`
	MusicTime   int    `json:"time"`
	IsLoveMusic bool   `json:"love"`
}

type MusicListPacket struct {
	MusicList []MusicElementPacket `json:"musicList"`
}

type MusicInfoPacket struct {
	ErrorCode int             `json:"code"`
	Param     MusicListPacket `json:"param"`
}

type LoveMusicElementPacket struct {
	MusicID    int    `json:"id"`
	MusicName  string `json:"musicname"`
	ArtistName string `json:"artist"`
	AlbumName  string `json:"albumname"`
	MusicTime  int    `json:"time"`
	LoveDegree int    `json:"degree"`
}

type LoveMusicListPacket struct {
	LoveMusicList []LoveMusicElementPacket `json:"musicList"`
}

type LoveMusicInfoPacket struct {
	ErrorCode int                 `json:"code"`
	Param     LoveMusicListPacket `json:"param"`
}

type MusicController struct {
}

func NewMusicController() *MusicController {
	return &MusicController{}
}

func (this *MusicController) readBodyByStream(body io.ReadCloser, w http.ResponseWriter) {
	reader := bufio.NewReader(body)
	content := make([]byte, HttpReadSize)
	totalSize := 0
	for {
		size, err := reader.Read(content)
		if size != 0 {
			totalSize += size
			w.Write(content[:size])
		} else if err != nil || size == 0 {
			break
		}
		w.(http.Flusher).Flush()
	}
}

func (this *MusicController) readBody(body io.ReadCloser, w http.ResponseWriter) {
	content, err := ioutil.ReadAll(body)
	if err == nil {
		w.Write(content)
	}
}

func (this *MusicController) fetchMusicWithProxy(url string, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	w.Header().Set("Connection", "keep-alive")
	defer resp.Body.Close()
	this.readBodyByStream(resp.Body, w)
}

func (this *MusicController) fetchRandomList(userId int, w http.ResponseWriter, r *http.Request) {
	channel, err := strconv.Atoi(r.Form.Get("channel"))
	if err != nil {
		NormalResponse(w, InvalidParam)
	} else {
		musicType, err := strconv.Atoi(r.Form.Get("type"))
		if err != nil || musicType > music.MusicTypeMax || musicType < 0 {
			NormalResponse(w, InvalidParam)
		} else {
			music.MusicManagerInstance().SetPlayer(music.BaiduPlayer)
			musicList, err := music.MusicManagerInstance().FetchMusicList(userId, channel, musicType)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				if len(musicList) != 0 {
					this.writeMusicInfo(musicList, w)
				}
			}
		}
	}
}

func (this *MusicController) writeMusicInfo(musicList []*element.MusicInfo, w http.ResponseWriter) {
	var packet *MusicInfoPacket = &MusicInfoPacket{}
	for _, music := range musicList {
		var elementPacket MusicElementPacket
		elementPacket.MusicID = music.MusicId
		elementPacket.MusicName = music.MusicName
		elementPacket.ArtistName = music.MusicAuthor
		elementPacket.AlbumName = music.AlbumName
		elementPacket.MusicTime = music.MusicTime
		elementPacket.IsLoveMusic = music.IsLoveMusic
		packet.Param.MusicList = append(packet.Param.MusicList, elementPacket)
		fmt.Println(music.MusicName)
	}
	body, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("fetch Random List Marshal Error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}

func (this *MusicController) fetchLoveList(userId int, w http.ResponseWriter, r *http.Request) {
	musicList, err := music.MusicManagerInstance().FetchLoveList(userId)
	// musicList, err := model.MusicModelInstance().FetchLoveList(userId)
	if err != nil {
		fmt.Println("fetchLoveList error: ", err)
		NormalResponse(w, DatabaseError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		this.writeMusicInfo(musicList, w)
	}
}

func (this *MusicController) fetchListenedList(userId int, w http.ResponseWriter, r *http.Request) {
	musicList, err := model.MusicModelInstance().FetchLoveList(userId)
	if err != nil {
		fmt.Println("fetchLoveList error: ", err)
		NormalResponse(w, DatabaseError)
	} else {
		var packet *LoveMusicInfoPacket = &LoveMusicInfoPacket{}
		for _, music := range musicList {
			var elementPacket LoveMusicElementPacket
			elementPacket.MusicID = music.MusicId
			elementPacket.MusicName = music.MusicName
			elementPacket.ArtistName = music.MusicAuthor
			elementPacket.AlbumName = music.AlbumName
			elementPacket.MusicTime = music.MusicTime
			elementPacket.LoveDegree = music.LoveDegree
			packet.Param.LoveMusicList = append(packet.Param.LoveMusicList, elementPacket)
			fmt.Println(music.MusicName)
		}
		body, err := json.Marshal(packet)
		if err != nil {
			fmt.Println("fetch Random List Marshal Error: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	}
}

func (this *MusicController) loveMusic(userId int, w http.ResponseWriter, r *http.Request) {
	musicId, err1 := strconv.Atoi(r.Form.Get("musicId"))
	loveDegree, err2 := strconv.Atoi(r.Form.Get("degree"))
	if err1 != nil || err2 != nil {
		NormalResponse(w, InvalidParam)
	} else {
		if model.MusicModelInstance().LoveMusic(userId, musicId, loveDegree) != nil {
			NormalResponse(w, DatabaseError)
		} else {
			NormalResponse(w, OK)
		}
	}
}

func (this *MusicController) listenMusic(userId int, w http.ResponseWriter, r *http.Request) {
	musicId, err := strconv.Atoi(r.Form.Get("musicId"))
	if err != nil {
		NormalResponse(w, InvalidParam)
	} else {
		if err = model.MusicModelInstance().ListenMusic(userId, musicId); err != nil {
			fmt.Println(err)
			NormalResponse(w, DatabaseError)
		} else {
			NormalResponse(w, OK)
		}
	}
}

func (this *MusicController) searchMusic(userId int, w http.ResponseWriter, r *http.Request) {
	key := r.Form.Get("key")
	fmt.Println("key = ", key)
	searchList, err := music.MusicManagerInstance().SearchMusic(userId, key)
	if err != nil {
		fmt.Println(err)
		NormalResponse(w, DatabaseError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		this.writeMusicInfo(searchList, w)
	}
}

func (this *MusicController) MusicAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("cookie Error: ", err)
		NormalResponse(w, InvalidToken)
	} else if len(tokenCookie.Value) == 0 {
		NormalResponse(w, EmptyToken)
	} else {
		findToken, userId := token.TokenManagerInstance().CheckTokenExist(tokenCookie.Value)
		if findToken == false {
			fmt.Println("search Token Error: ", err)
			NormalResponse(w, InvalidToken)
		} else {
			fmt.Println("userId: ", userId)
			action := r.Form.Get("action")
			fmt.Println(action)
			switch action {
			case "fetchRandomList":
				this.fetchRandomList(userId, w, r)
			case "fetchLoveList":
				this.fetchLoveList(userId, w, r)
			case "fetchListenedList":
				this.fetchListenedList(userId, w, r)
			case "loveMusic":
				fmt.Println("loveMusic")
				this.loveMusic(userId, w, r)
			case "listenMusic":
				this.listenMusic(userId, w, r)
			case "searchMusic":
				this.searchMusic(userId, w, r)
			case "test":
				player.QQMusicPlayerInstance().SearchMusic("周杰伦", 0, 50)
			case "musicProxy":
				// 先测试请求某一首歌
				// action="musicProxy&player=baidu&type=fetchList&"
				var musicURL string = "http://yinyueshiting.baidu.com/data2/music/134380410/5963228158400128.mp3?xcode=55a3d05c2c746155f1962d1d1690fc93"
				this.fetchMusicWithProxy(musicURL, w, r)
			}
		}
	}
}
