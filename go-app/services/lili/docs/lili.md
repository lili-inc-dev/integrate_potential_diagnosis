### 1. N/A

1. route definition

- Url: /admin
- Method: POST
- Request: `CreateAdminReq`
- Response: `CreateAdminRes`

2. request definition



```golang
type CreateAdminReq struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```


3. response definition



```golang
type CreateAdminRes struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```

### 2. N/A

1. route definition

- Url: /admin
- Method: PUT
- Request: `UpdateAdminReq`
- Response: `UpdateAdminRes`

2. request definition



```golang
type UpdateAdminReq struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```


3. response definition



```golang
type UpdateAdminRes struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```

### 3. N/A

1. route definition

- Url: /admin
- Method: DELETE
- Request: `DeleteAdminReq`
- Response: `DeleteAdminRes`

2. request definition



```golang
type DeleteAdminReq struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type DeleteAdminRes struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```

### 4. N/A

1. route definition

- Url: /admin
- Method: GET
- Request: `FindByIdAdminReq`
- Response: `FindByIdAdminRes`

2. request definition



```golang
type FindByIdAdminReq struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type FindByIdAdminRes struct {
	Id uint64 `json:"id"`
	RoleId int64 `json:"role_id"`
	FirebaseUid string `json:"firebase_uid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Affiliation string `json:"affiliation"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IconUrl string `json:"icon_url"`
	Description string `json:"description"`
}
```

