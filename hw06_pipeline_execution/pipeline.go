package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageIn := in
	var stageOut Out
	for _, stage := range stages {
		stageOut = stage(stageIn)
		stageIn = stageOut
	}
	summary := make(Bi)
	go func() {
		defer close(summary)
		for {
			select {
			case <-done:
				return
			case i, ok := <-stageOut:
				if ok {
					summary <- i
				} else {
					return
				}
			}
		}
	}()
	return summary
}
