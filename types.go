package pkg

type PayloadAction = string

var (
	ActionSubmit        PayloadAction = "submit"
	ActionCancel        PayloadAction = "cancel"
	ActionCancelWithKey PayloadAction = "cancel_with_key" // this method will not notify to the user
	ActionReload        PayloadAction = "reload"
)
