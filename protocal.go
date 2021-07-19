package main

type Ret struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type List struct {
	List interface{} `json:"list"`
}
