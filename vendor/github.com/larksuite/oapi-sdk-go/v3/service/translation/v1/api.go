// Package translation code generated by oapi sdk gen
/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larktranslation

import (
	"context"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func NewService(config *larkcore.Config) *TranslationService {
	t := &TranslationService{config: config}
	t.Text = &text{service: t}
	return t
}

type TranslationService struct {
	config *larkcore.Config
	Text   *text // 文本
}

type text struct {
	service *TranslationService
}

// 文本语种识别
//
// - 机器翻译 (MT)，支持 100 多种语言识别，返回符合 ISO 639-1 标准
//
// - 单租户限流：20QPS，同租户下的应用没有限流，共享本租户的 20QPS 限流
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/ai/translation-v1/text/detect
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/translationv1/detect_text.go
func (t *text) Detect(ctx context.Context, req *DetectTextReq, options ...larkcore.RequestOptionFunc) (*DetectTextResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/translation/v1/text/detect"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, t.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &DetectTextResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, t.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 文本翻译
//
// - 机器翻译 (MT)，支持以下语种互译：;"zh": 汉语；;"zh-Hant": 繁体汉语；;"en": 英语；;"ja": 日语；;"ru": 俄语；;"de": 德语；;"fr": 法语；;"it": 意大利语；;"pl": 波兰语；;"th": 泰语；;"hi": 印地语；;"id": 印尼语；;"es": 西班牙语；;"pt": 葡萄牙语；;"ko": 朝鲜语；;"vi": 越南语；
//
// - 单租户限流：20QPS，同租户下的应用没有限流，共享本租户的 20QPS 限流
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/ai/translation-v1/text/translate
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/translationv1/translate_text.go
func (t *text) Translate(ctx context.Context, req *TranslateTextReq, options ...larkcore.RequestOptionFunc) (*TranslateTextResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/translation/v1/text/translate"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, t.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &TranslateTextResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, t.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
