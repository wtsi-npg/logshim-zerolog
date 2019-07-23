/*
 * Copyright (C) 2019. Genome Research Ltd. All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @file log.go
 * @author Keith James <kdj@sanger.ac.uk>
 */

package logshim_zerolog

import (
	"fmt"
	"io"
	"time"

	"github.com/kjsanger/logshim"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func translateLevel(level logshim.Level) (zerolog.Level, error) {
	var (
		lvl zerolog.Level
		err error
	)

	switch level {
	case logshim.ErrorLevel:
		lvl = zerolog.ErrorLevel
	case logshim.WarnLevel:
		lvl = zerolog.WarnLevel
	case logshim.NoticeLevel:
		fallthrough
	case logshim.InfoLevel:
		lvl = zerolog.InfoLevel
	case logshim.DebugLevel:
		lvl = zerolog.DebugLevel
	default:
		lvl = zerolog.WarnLevel
		err = errors.New(fmt.Sprintf("invalid log level %d, defaulting to WARN level", level))
	}

	return lvl, err
}

type ZeroLogger struct {
	name string
	*zerolog.Logger
}

func New(writer io.Writer, level logshim.Level) *ZeroLogger {
	lvl, err := translateLevel(level)
	lg := zerolog.New(writer).Level(lvl).With().Timestamp().Logger()

	if err != nil {
		lg.Error().Err(err).Msg("log configuration error")
	}
	return &ZeroLogger{"ZeroLog", &lg}
}

func (log *ZeroLogger) Name() string {
	return log.name
}

func (log *ZeroLogger) Err(err error) logshim.Message {
	return &zeroMessage{log.Logger.Err(err)}
}

func (log *ZeroLogger) Error() logshim.Message {
	return &zeroMessage{log.Logger.Error()}
}

func (log *ZeroLogger) Warn() logshim.Message {
	return &zeroMessage{log.Logger.Warn()}
}

func (log *ZeroLogger) Debug() logshim.Message {
	return &zeroMessage{log.Logger.Debug()}
}

func (log *ZeroLogger) Notice() logshim.Message {
	return &zeroMessage{log.Logger.Info()}
}

func (log *ZeroLogger) Info() logshim.Message {
	return &zeroMessage{log.Logger.Info()}
}

type zeroMessage struct {
	*zerolog.Event
}

func (msg *zeroMessage) Err(err error) logshim.Message {
	msg.Event.Err(err)
	return msg
}

func (msg *zeroMessage) Bool(key string, val bool) logshim.Message {
	msg.Event.Bool(key, val)
	return msg
}

func (msg *zeroMessage) Int(key string, val int) logshim.Message {
	msg.Event.Int(key, val)
	return msg
}

func (msg *zeroMessage) Int64(key string, val int64) logshim.Message {
	msg.Event.Int64(key, val)
	return msg
}

func (msg *zeroMessage) Uint64(key string, val uint64) logshim.Message {
	msg.Event.Uint64(key, val)
	return msg
}

func (msg *zeroMessage) Str(key string, val string) logshim.Message {
	msg.Event.Str(key, val)
	return msg
}

func (msg *zeroMessage) Time(key string, val time.Time) logshim.Message {
	msg.Event.Time(key, val)
	return msg
}

func (msg *zeroMessage) Msg(val string) {
	msg.Event.Msg(val)
}

func (msg *zeroMessage) Msgf(format string, a ...interface{}) {
	msg.Msg(fmt.Sprintf(format, a...))
}
