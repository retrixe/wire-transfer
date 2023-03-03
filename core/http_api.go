package core

type WsUploadRequest struct {
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Hash      string `json:"hash"`
	PublicKey string `json:"public_key"`
	Port      int    `json:"port"`
	Precache  bool   `json:"precache"`
}

type WsUploadResponse struct {
	Token    string `json:"token"`
	Precache bool   `json:"precache"`
}

type WsUploadProxiedTransferRequest struct {
	ProxyToken string `json:"proxy_token"`
}

type GetInfoResponse struct {
	Name            string `json:"name"`
	Size            int    `json:"size"`
	Hash            string `json:"hash"`
	Available       bool   `json:"available"`
	SupportsDirect  bool   `json:"supports_direct"`
	SupportsProxied bool   `json:"supports_proxied"`
}

type GetDownloadDirectResponse struct {
	Host string `json:"host"`
	Size int    `json:"size"`
	Hash string `json:"hash"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type GetDownloadProxiedResponse struct {
	Host  string `json:"host"`
	Size  int    `json:"size"`
	Hash  string `json:"hash"`
	Token string `json:"token"`
}
