package controller

type ContextDummyFlashStorage struct {
	errorBag *ErrorBag
	inputBag *InputBag
}

func NewContextDummyFlashStorage() ContextFlashStorage {
	return &ContextDummyFlashStorage{
		errorBag: newErrorBag(),
		inputBag: newInputBag(),
	}
}

func (c *ContextDummyFlashStorage) Flush() {
}

func (c *ContextDummyFlashStorage) Errors() *ErrorBag {
	return c.errorBag
}

func (c *ContextDummyFlashStorage) Inputs() *InputBag {
	return c.inputBag
}
