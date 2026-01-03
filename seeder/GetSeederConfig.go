package seeder

import "github.com/AlecAivazis/survey/v2"

type SeederOptions struct {
	InsertData        bool
	InsertNulls       bool
	InsertOutOfBounds bool
	SamplesCount      int
}

func GetSeederConfig() (SeederOptions, error) {
	var opts SeederOptions

	if err := survey.AskOne(
		&survey.Confirm{Message: "Insert example data?"},
		&opts.InsertData,
	); err != nil {
		return opts, err
	}

	if !opts.InsertData {
		return opts, nil
	}

	survey.AskOne(
		&survey.Confirm{Message: "Insert NULLs to emulate sensor faults?"},
		&opts.InsertNulls,
	)

	survey.AskOne(
		&survey.Confirm{Message: "Insert values in warning/danger levels?"},
		&opts.InsertOutOfBounds,
	)

	survey.AskOne(
		&survey.Input{Message: "Samples count:"},
		&opts.SamplesCount,
	)

	return opts, nil
}
