package entities

type ErrDevEnvNotExists struct {
	ClusterName string
	DevEnvName  string
}

func (ErrDevEnvNotExists) Error() string {
	return "ErrDevEnvNotExists"
}

type ErrStartRemovingDevEnv struct {
	DevEnvName string
}

func (ErrStartRemovingDevEnv) Error() string {
	return "ErrStartRemovingDevEnv"
}

type ErrStartStoppingDevEnv struct {
	DevEnvName string
}

func (ErrStartStoppingDevEnv) Error() string {
	return "ErrStartStoppingDevEnv"
}

type ErrStopRemovingDevEnv struct {
	DevEnvName string
}

func (ErrStopRemovingDevEnv) Error() string {
	return "ErrStopRemovingDevEnv"
}

type ErrStopCreatingDevEnv struct {
	DevEnvName string
}

func (ErrStopCreatingDevEnv) Error() string {
	return "ErrStopCreatingDevEnv"
}

type ErrStopStartingDevEnv struct {
	DevEnvName string
}

func (ErrStopStartingDevEnv) Error() string {
	return "ErrStopStartingDevEnv"
}
