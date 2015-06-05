package limit

import (
	"github.com/AdOnWeb/postmanq/common"
	"github.com/AdOnWeb/postmanq/log"
	"sync/atomic"
)

type Limiter struct {
	id int
}

func (l *Limiter) run(id int) {
	l.id = id
	for event := range events {
		l.check(event)
	}
}

func (l *Limiter) check(event *common.SendEvent) {
	log.Info("limiter#%d check limit for mail#%d", l.id, event.Message.Id)
	// пытаемся найти ограничения для почтового сервиса
	if limit, ok := service.Limits[event.Message.HostnameTo]; ok {
		log.Debug("limiter#%d found config for %s", l.id, event.Message.HostnameTo)
		// если оно нашлось, проверяем, что отправка нового письма происходит в тот промежуток времени,
		// в который нам необходимо следить за ограничениями
		if limit.isValidDuration(event.Message.CreatedDate) {
			atomic.AddInt32(&limit.currentValue, 1)
			currentValue := atomic.LoadInt32(&limit.currentValue)
			log.Debug("limiter#%d get current value %d, const value %d", l.id, currentValue, limit.Value)
			// если ограничение превышено
			if currentValue > limit.Value {
				log.Debug("limiter#%d current value is exceeded for %s", l.id, event.Message.HostnameTo)
				// определяем очередь, в которое переложем письмо
				event.Message.BindingType = limit.bindingType
				// говорим получателю, что у нас превышение ограничения,
				// разблокируем поток получателя
				event.Result <- OverlimitSendEventResult
				return
			}
		} else {
			log.Debug("limiter#%d duration great then %v", l.id, limit.duration)
		}
	} else {
		log.Debug("limiter#%d not found for %s", l.id, event.Message.HostnameTo)
	}
	event.Iterator.Next().(common.SendingService).Events() <- event
}
