package runner

import (
	"github.com/Khorevaa/go-v8runner/find"
	"github.com/Khorevaa/go-v8runner/types"
	"os/exec"
	"strconv"
	"strings"

	"context"
	"github.com/Khorevaa/go-v8runner/errors"
	"time"
)

var VERSION_1S = "8.3"

type Runner struct {
	Options *Options
}

func NewRunner() *Runner {

	return &Runner{
		Options: defaultOptions(),
	}

}

func (r *Runner) Run(where types.InfoBase, what types.Command, opts ...interface{}) error {

	options := r.Options
	inlineOptions := getOptions(opts...)
	if inlineOptions != nil {
		options = inlineOptions
	}

	if options == nil {
		options = defaultOptions()
	}

	o := clearOpts(opts)

	options.Options(o...)

	return r.RunWithOptions(where, what, options)

}

func (r *Runner) RunWithOptions(where types.InfoBase, what types.Command, options *Options) error {

	checkErr := what.Check()
	if checkErr != nil {
		return checkErr
	}

	commandV8, err := getV8Path(options)

	if err != nil {
		return err
	}

	values := make(types.Values)

	connectString := where.Values()
	values.Set("/IBConnectionString", types.SpaceSep, strings.Join(connectString.Values(), ";"))

	if what.Command() == types.COMMAND_CREATEINFOBASE {
		connectString.Append(what.Values())
	} else {
		values.Append(what.Values())
	}

	values.Append(options.commonValues)
	values.Append(getOptionsValues(options))
	values.Append(options.customValues)

	var args []string

	args = append(args, what.Command())
	if what.Command() == types.COMMAND_CREATEINFOBASE {

		values.Del("/IBConnectionString")
		args = append(args, strings.Join(connectString.Values(), ";"))

		clearValuesForCreateInfobase(&values)
	}
	args = append(args, values.Values()...)

	defer options.RemoveTempFiles()

	var errRun error

	if options.Context != nil {
		errRun = runCommandContext(options.Context, commandV8, args, options)
	} else {
		errRun = runCommand(commandV8, args, options)
	}

	if errRun != nil {
		outLog, _ := readOut(options.Out)
		errorWithOut := errors.AddErrorContext(errRun, "out", outLog)
		return errorWithOut
	}

	dumpErr := checkRunResult(options)

	if dumpErr != nil {
		return dumpErr
	}

	return nil
}

func getV8Path(options *Options) (string, error) {
	if len(options.v8path) > 0 {
		return options.v8path, nil
	}

	v8 := VERSION_1S
	if len(options.Version) > 0 {
		v8 = options.Version
	}

	v8path, err := find.PlatformByVersion(v8, find.WithBitness(find.V8_x64x32))

	if err != nil {

		err = errors.NotExist.Newf("Version %s not found", options.Version)
		_ = errors.AddErrorContext(err, "version", options.Version)

		return "", err
	}

	return v8path, nil

}

func readOut(file string) (string, error) {
	b, err := readV8file(file)
	return string(b), err
}

func readDumpResult(file string) int {

	b, _ := readV8file(file)
	code, _ := strconv.ParseInt(string(b), 10, 64)
	return int(code)
}

func runCommand(command string, args []string, opts *Options) (err error) {

	cmd := exec.Command(command, args...)
	err = cmd.Run()

	if err != nil {
		errors.Runtime.Wrapf(err, "run command exec error")
	}

	return
}

func runCommandContext(ctx context.Context, command string, args []string, opts *Options) (err error) {

	runCtx := ctx
	if opts.Timeout > 0 {

		timeout := int64(time.Second) * opts.Timeout
		ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout))
		defer cancel() // The cancel should be deferred so resources are cleaned up
		runCtx = ctxTimeout
	}

	cmd := exec.CommandContext(runCtx, command, args...)

	err = cmd.Run()

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		err = errors.Timeout.Wrap(err, "run command context timeout exceeded")
	case err != nil:
		err = errors.Runtime.Wrap(err, "run command context exec error")
	}

	return
}

func checkRunResult(options *Options) error {

	dumpCode := readDumpResult(options.DumpResult)

	switch dumpCode {

	case 0:

		return nil

	case 1:

		outLog, _ := readOut(options.Out)
		err := errors.Internal.New("1S internal error")
		errorWithOut := errors.AddErrorContext(err, "out", outLog)
		return errorWithOut

	case 101:

		outLog, _ := readOut(options.Out)
		err := errors.Database.New("1S database error")
		errorWithOut := errors.AddErrorContext(err, "out", outLog)
		return errorWithOut

	default:

		outLog, _ := readOut(options.Out)
		err := errors.Invalid.New("Unknown 1S error")
		errorWithOut := errors.AddErrorContext(err, "out", outLog)
		return errorWithOut

	}

}

func clearValuesForCreateInfobase(v *types.Values) {

	// TODO Сделать очистку значение

}

func getOptionsValues(options *Options) types.Values {

	values := make(types.Values)

	outValue := options.Out
	if options.NoTruncate {
		outValue += " -NoTruncate"
	}

	values.Set("/Out", types.SpaceSep, outValue)
	values.Set("/DumpResult", types.SpaceSep, options.DumpResult)

	return values

}

func defaultOptions() *Options {

	options := Options{}

	options.NewOutFile()
	options.NewDumpResultFile()
	options.customValues = make(types.Values)
	options.commonValues = make(types.Values)

	return &options
}

func getOptions(opts ...interface{}) *Options {

	for _, opt := range opts {

		switch opt.(type) {

		case Options:
			userOptions, _ := opt.(Options)
			return &userOptions
		case *Options:
			userOptions, _ := opt.(*Options)
			return userOptions
		}

	}

	return nil
}

func clearOpts(opts []interface{}) []Option {

	var o []Option

	for _, opt := range opts {

		if fn, ok := opt.(Option); ok {
			o = append(o, fn)
		}
	}
	return o
}