package models

type EditorFile struct {
	EditorFileId int32  `db:"editor_file_id"`
	AccountId    int32  `db:"account_id"`
	FileName     string `db:"file_name"`
	FileContent  string `db:"file_content"`
}
