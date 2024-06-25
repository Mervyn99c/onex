// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package sms

import (
	"context"
	"encoding/json"
	"fmt"
	message2 "github.com/Rosas99/smsx/internal/sms/biz/message"
	"github.com/Rosas99/smsx/internal/sms/model"
	"github.com/Rosas99/smsx/internal/sms/store"
	"github.com/segmentio/kafka-go"
	"time"

	"github.com/Rosas99/smsx/pkg/log"
	genericoptions "github.com/Rosas99/smsx/pkg/options"
)

// kafkaLogger is a log.Logger implementation that writes log messages to Kafka.
type Logger struct {
	// enabled is an atomic boolean indicating whether the logger is enabled.
	enabled int32
	// writer is the Kafka writer used to write log messages.

	writer *kafka.Writer
	ds     store.TemplateStore
}

// todo 这里改成短信发送历史
// AuditMessage is the message structure for log messages.
type AuditMessage struct {
	Matcher   string     `protobuf:"bytes,1,opt,name=matcher,proto3" json:"matcher,omitempty"`
	Request   []any      `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	Result    bool       `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	Explains  [][]string `protobuf:"bytes,4,opt,name=explains,proto3" json:"explains,omitempty"`
	Timestamp int64      `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// NewLogger creates a new kafkaLogger instance.
func NewLogger(kafkaOpts *genericoptions.KafkaOptions, ds store.TemplateStore) (*Logger, error) {
	writer, err := kafkaOpts.Writer()
	if err != nil {
		return nil, err
	}

	return &Logger{writer: writer, ds: ds}, nil
}

// LogModel writes a log message for the policy model.
func (l *Logger) LogHistory(test string) string {

	//message := AuditMessage{
	//	Timestamp: time.Now().Unix(),
	//}

	// 这里先转成json
	out, _ := json.Marshal(test)
	// 因为writer传的都是二进制
	// 这里json转成二进制数组
	//fmt.Println(test)
	//if err := l.writer.WriteMessages(context.Background(), kafka.Message{Value: out}); err != nil {
	//	log.Errorw(err, "Failed to write kafka messages")
	//} else {
	//	fmt.Println(string(out))
	//}

	err := l.ds.Create(context.Background(), &model.TemplateM{})
	if err != nil {

	}

	return string(out)
	// todo 修改成MySQL
}

// log others

// LogModel writes a log message for the policy model.
func (l *Logger) LogMessageRequest(test string) string {

	message := message2.TemplateMsgRequest{
		Timestamp: time.Now().Unix(),
	}

	// 这里先转成json
	out, _ := json.Marshal(message)
	// 因为writer传的都是二进制
	// 这里json转成二进制数组
	fmt.Println(test)
	if err := l.writer.WriteMessages(context.Background(), kafka.Message{Value: out}); err != nil {
		log.Errorw(err, "Failed to write kafka messages")
	} else {
		fmt.Println(string(out))
	}
	return string(out)
}
