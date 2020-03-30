package save

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	ConfirmMessage  = "Would you like to\nsave the game?"
	SavingMessage   = "SAVING. DON'T TURN\nOFF THE POWER."
	CompleteMessage = "%s saved\nthe game."
)

// Data - セーブデータ情報を格納する
type Data struct {
	Point    Point    `json:"point"`
	Avatar   Avatar   `json:"avatar"`
	FlagData FlagData `json:"flag"`
	Valid    bool     // セーブデータが有効か
}

// Point - ユーザーがセーブをした場所
type Point struct {
	Stage string `json:"stage"`
	Index int    `json:"index"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

// Avatar - ユーザーのアバターデータ
type Avatar struct {
	ID     int    `json:"id"`     // アバターID
	Name   string `json:"name"`   // アバターネーム
	Minute uint   `json:"minute"` // プレイ時間(分)
}

// FlagData - ユーザーのフラグデータ
type FlagData struct {
}

// New - コンストラクタ
func New(filename string) *Data {
	savedata := &Data{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return CreateNewData(1, "ethan")
	}

	if err := json.Unmarshal(file, savedata); err != nil {
		return CreateNewData(1, "ethan")
	}

	savedata.Valid = true

	return savedata
}

// Write - SaveDataを保存する
func Write(sav *Data) {
	jsonData, _ := json.MarshalIndent(sav, "", "    ")
	ioutil.WriteFile("savedata.json", jsonData, os.ModePerm)
}

// CreateNewData - 新しいセーブデータを作る
func CreateNewData(avatarID int, avatarName string) *Data {
	savedata := &Data{
		Point: Point{
			Stage: "Oxalis City",
			Index: 1,
			X:     37,
			Y:     16,
		},
		Avatar: Avatar{
			ID:     avatarID,
			Name:   avatarName,
			Minute: 0,
		},
		FlagData: FlagData{},
	}
	Write(savedata)
	return savedata
}

// Message - セーブメッセージを返す
func Message(name string) []string {
	return []string{
		ConfirmMessage,
		SavingMessage,
		fmt.Sprintf(CompleteMessage, name),
		":end",
	}
}
