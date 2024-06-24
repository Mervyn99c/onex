// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package sms // import "github.com/superproj/onex/pkg/fakeserver"

// validate 暂时不用

// todo 业务消息处理放sms
// 日志采集放pump

// 使用wire

// 参考user center 只使用jwt鉴权 作为admin接口 需要登录才能使用创建和配置模板
// 发送短信不用
// jwt密钥如何保存和生成？

// kpi日志如何实现可视化

// 先使用dummy供应商
// 对接阿里云 腾讯云短信
// 对接短信报告 短信报告的用处
// 调研批量接口

// 使用gin中间件生成幂等id 将id传到mq mq进行check

// mq先new redis kafka 然后logic使用结构体组装redis等

// 了解context的使用流程 // 目前都没用到context的key value

// mq考虑先使用kquene 后续再考虑如何集成pump
// 从消息获取幂等id 检查幂等id
// 发送短信：获取供应商
// 主备供应商 有一个发送成功即停止 不然下一个
// 短信发送历史

// retry 调用接口重试

// 规则使用导包的方式
// todo 接入usercenter pump watcher
//usercenter 暂时只使用认证功能，授权功能后续再使用

// 使用的是golang-jwt包，后续考虑如何集成gin-jwt

// 了解装饰器的洋葱模型 VS gin的pipeline模式
// https://jianghushinian.cn/2022/02/28/Golang-%E5%B8%B8%E8%A7%81%E8%AE%BE%E8%AE%A1%E6%A8%A1%E5%BC%8F%E4%B9%8B%E8%A3%85%E9%A5%B0%E6%A8%A1%E5%BC%8F/

// 参考gateway的中间件，将usercenter认证作为rpc调用
// 登录等其他接口直接对usercenter进行http调用

// 熟悉gin的各种使用

// 最终考虑
