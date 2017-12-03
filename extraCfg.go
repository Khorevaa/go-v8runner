package v8runner

import "fmt"

func (conf *Конфигуратор) ЗапуститьВРежимеПредприятияСКлючемЗапуска(КлючЗапуска string, УправляемыйРежим bool, ДополнительныеПараметры ...string) (err error) {

	ДополнительныеПараметры = append(ДополнительныеПараметры, fmt.Sprintf("/C%s", КлючЗапуска))

	err = conf.ЗапуститьВРежимеПредприятия(УправляемыйРежим, ДополнительныеПараметры...)

	return
}

func (conf *Конфигуратор) ЗапуститьВРежимеПредприятия(УправляемыйРежим bool, ДополнительныеПараметры ...string) (err error) {

	var Параметры []string

	if УправляемыйРежим {
		Параметры = append(Параметры, "/RunModeManagedApplication")
	} else {
		Параметры = append(Параметры, "/RunModeOrdinaryApplication")
	}
	Параметры = append(Параметры, ДополнительныеПараметры...)

	conf.УстановитьПараметры(Параметры...)
	err = conf.ВыполнитьКомандуПредприятие()

	return
}

//LoadExternalDataProcessorOrReportFromFiles
func (conf *Конфигуратор) СобратьОбработкуОтчетИзФайлов(ПапкаИсходников string, ИмяФайлаОбработки string, ДополнительныеПараметры ...string) (err error) {

	var Параметры []string

	Параметры = append(Параметры, "/LoadExternalDataProcessorOrReportFromFiles", ПапкаИсходников, ИмяФайлаОбработки)
	Параметры = append(Параметры, ДополнительныеПараметры...)

	conf.УстановитьПараметры(Параметры...)
	err = conf.ВыполнитьКоманду()

	return
}
