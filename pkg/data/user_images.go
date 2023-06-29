package data

type UserImage struct {
    ID int `json:"id"`
    UserID string `json:"user_id"`
    FileName string `json:"file_name"`
    CreatedAt string `json:"-"`
    UpdatedAt string `json:"-"`
}
