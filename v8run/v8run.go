package v8run

import (
	"fmt"
	"os/exec"
	"syscall"

	"../v8platform"
	"../v8tools"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//noinspection NonAsciiCharacters
type ЗапускательКонфигуратора struct {
	файлИнформации                   string
	очищатьФайлИнформации            bool
	этоWindows                       bool
	версияПлатформы                  *v8platform.ВерсияПлатформы
	ключСоединенияСБазой             string
	пользовательскиеПараметрыЗапуска []string
	параметыЗапуска                  []string
	параметрыАвторизации             *параметрыАвторизации
	командаКонфигуратора             командаКонфигуратора
	выводКоманды                     string
}

const (
	// Типы протоколов подключения

	КомандаКонфигуратор = командаКонфигуратора("DESIGNER")
	КомандаСоздатьБазу  = командаКонфигуратора("CREATEINFOBASE")
	КомандаПредприятие  = командаКонфигуратора("ENTERPRISE")

	КомандаПоУмолчанию = КомандаКонфигуратор
)

type командаКонфигуратора string

var доступныеКомандыКонфигуратора = []командаКонфигуратора{КомандаКонфигуратор, КомандаСоздатьБазу, КомандаПредприятие}

type параметрыАвторизации struct {
	Пользователь string
	Пароль       string
}

func (conf *ЗапускательКонфигуратора) ВыполнитьКомандуКонфигуратора() (err error) {

	conf.командаКонфигуратора = КомандаКонфигуратор
	err = conf.запуститьКоманду()
	return
}

func (conf *ЗапускательКонфигуратора) ВыполнитьКомандуСоздатьБазу() (err error) {

	conf.командаКонфигуратора = КомандаСоздатьБазу
	err = conf.запуститьКоманду()
	return
}

func (conf *ЗапускательКонфигуратора) ВыполнитьКомандуПредприятие() (err error) {

	conf.командаКонфигуратора = КомандаПредприятие
	err = conf.запуститьКоманду()
	return
}

func (conf *ЗапускательКонфигуратора) ВыполнитьКоманду() (err error) {

	err = conf.ВыполнитьКомандуКонфигуратора()
	return
}

func (conf *ЗапускательКонфигуратора) запуститьКоманду() (err error) {

	conf.собратьПараметрыЗапуска()

	_, checkErr := conf.ПроверитьВозможностьВыполнения()

	if checkErr != nil {
		return
	}

	ok, err := conf.проверитьВозможностьВыполнения()

	if ok {
		err = conf.выполнить(conf.параметыЗапуска)
	}
	return
}

//export func

func (conf *ЗапускательКонфигуратора) УстановитьВерсиюПлатформы(строкаВерсияПлатформы string) {

	conf.версияПлатформы = v8platform.ПолучитьВерсию(строкаВерсияПлатформы)

}

func (conf *ЗапускательКонфигуратора) КлючСоединенияСБазой() string {

	log.Debugf("Получение ключа соединения с базой: %s", conf.ключСоединенияСБазой)

	return conf.ключСоединенияСБазой
}

func (conf *ЗапускательКонфигуратора) УстановитьКлючСоединенияСБазой(КлючСоединенияСБазой string) {

	log.Debugf("Установка ключа соединения с базой: %s", КлючСоединенияСБазой)

	conf.ключСоединенияСБазой = КлючСоединенияСБазой

}

func (conf *ЗапускательКонфигуратора) УстановитьАвторизацию(Пользователь string, Пароль string) {

	if conf.параметрыАвторизации == nil {
		conf.параметрыАвторизации = &параметрыАвторизации{}
	}

	conf.параметрыАвторизации.Пользователь = Пользователь
	conf.параметрыАвторизации.Пользователь = Пароль
}

func (conf *ЗапускательКонфигуратора) УстановитьПараметры(Параметры ...string) {

	conf.пользовательскиеПараметрыЗапуска = Параметры

}

func (conf *ЗапускательКонфигуратора) ДобавитьПараметры(Параметры ...string) {

	conf.пользовательскиеПараметрыЗапуска = append(conf.пользовательскиеПараметрыЗапуска, Параметры...)

}
func (c *ЗапускательКонфигуратора) ПолучитьВыводКоманды() (s string) {
	return c.выводКоманды
}

func (c *ЗапускательКонфигуратора) ПроверитьВозможностьВыполнения() (ok bool, err error) {

	return
}

func (conf *ЗапускательКонфигуратора) добавитьВыводВФайл() {

	if len(conf.файлИнформации) == 0 {
		conf.файлИнформации = v8tools.НовыйФайлИнформации()
	}

	conf.параметыЗапуска = append(conf.параметыЗапуска, "/Out", conf.файлИнформации)

	if !conf.очищатьФайлИнформации {
		conf.параметыЗапуска = append(conf.параметыЗапуска, "-NoTruncate")
	}

}
func (conf *ЗапускательКонфигуратора) добавитьАвторизацию() {

	Авторизации := conf.параметрыАвторизации

	if Авторизации == nil {
		return
	}

	if v8tools.ЗначениеЗаполнено(Авторизации.Пользователь) {
		conf.параметыЗапуска = append(conf.параметыЗапуска, fmt.Sprintf("/N %s", Авторизации.Пользователь))
	}
	if v8tools.ЗначениеЗаполнено(Авторизации.Пароль) {
		conf.параметыЗапуска = append(conf.параметыЗапуска, fmt.Sprintf("/P %s", Авторизации.Пароль))
	}

}

func (conf *ЗапускательКонфигуратора) собратьПараметрыЗапуска() {

	//conf.параметыЗапуска

	conf.параметыЗапуска = append(conf.параметыЗапуска, string(conf.командаКонфигуратора))

	if conf.командаКонфигуратора == КомандаСоздатьБазу {
		// TODO Сделать замену /F на File= или /S на Server=
	} else {
		conf.параметыЗапуска = append(conf.параметыЗапуска, conf.КлючСоединенияСБазой())
	}

	conf.добавитьАвторизацию()

	conf.параметыЗапуска = append(conf.параметыЗапуска, "/DisableStartupMessages")
	conf.параметыЗапуска = append(conf.параметыЗапуска, "/DisableStartupDialogs")

	conf.параметыЗапуска = append(conf.параметыЗапуска, conf.пользовательскиеПараметрыЗапуска...)

}

// private run func
const defaultFailedCode = 1

func (conf *ЗапускательКонфигуратора) выполнить(args []string) (e error) {

	var exitCode int

	procName := conf.версияПлатформы.V8
	cmd := exec.Command(procName, args...) // strings.Join(args, " "))

	log.Debugf("Строка запуска: %s", cmd.Args)

	out, e := cmd.Output()

	if e != nil {
		// try to get the exit code
		if exitError, ok := e.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Debugf("Could not get exit code for failed program: %v, %v", procName, args)
			exitCode = defaultFailedCode
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	conf.установитьВыводКоманды(conf.прочитатьФайлИнформации())

	if exitCode != 0 {
		e = errors.New(conf.выводКоманды)
	}

	log.Debugf("КодЗавершения команды: %v", exitCode)
	log.Debugf("Результат выполнения команды: %s, out: %s", conf.выводКоманды, out)
	return e

}

func (c *ЗапускательКонфигуратора) проверитьВозможностьВыполнения() (ok bool, err error) {

	ok = true

	if c.версияПлатформы == nil {
		err = errors.Wrap(err, "Не найдена доступная версия платформы")
		ok = false
	}

	return

}

func (c *ЗапускательКонфигуратора) установитьВыводКоманды(s string) {
	c.выводКоманды = s
	log.Debugf("Установлен вывод команды 1С: %s", s)
}

func (c *ЗапускательКонфигуратора) прочитатьФайлИнформации() (str string) {

	log.Debugf("Читаю файла информации 1С: %s", c.файлИнформации)

	b, err := v8tools.ReadFileUTF16(c.файлИнформации) // just pass the file name
	if err != nil {
		log.Debugf("Обшибка чтения файла информации 1С %s: %v", c.файлИнформации, err)
		str = ""
		return
		//fmt.Print(err)
	}

	str = string(b) // convert content to a 'string'

	return
}
