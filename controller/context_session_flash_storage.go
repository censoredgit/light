package controller

import "github.com/censoredgit/light/session"

type ContextSessionFlashStorage struct {
	sessionData *session.Data
	errorBag    *ErrorBag
	inputBag    *InputBag
}

func NewContextSessionFlashStorage(sessionData *session.Data) ContextFlashStorage {
	c := &ContextSessionFlashStorage{
		sessionData: sessionData,
		errorBag:    newErrorBag(),
		inputBag:    newInputBag(),
	}

	c.init()

	return c
}

func (c *ContextSessionFlashStorage) init() {
	if c.sessionData.Has("_error_input") {
		val := c.sessionData.Get("_error_input")
		if val != "" {
			errorInput, err := base64ToErrorBag(val)
			if err != nil {
				iLogError(err.Error())
			} else {
				c.errorBag.err = errorInput.Errors()
			}
			c.sessionData.Delete("_error_input")
		}
	}

	if c.sessionData.Has("_old_input") {
		val := c.sessionData.Get("_old_input")
		if val != "" {
			values, err := base64ToValues(val)
			if err != nil {
				iLogError(err.Error())
			} else {
				c.inputBag.old = values
			}
			c.sessionData.Delete("_old_input")
		}
	}

	if c.sessionData.Has("_input") {
		val := c.sessionData.Get("_input")
		if val != "" {
			values, err := base64ToValues(val)
			if err != nil {
				iLogError(err.Error())
			} else {
				c.inputBag.data = values
			}
			c.sessionData.Delete("_input")
		}
	}

}

func (c *ContextSessionFlashStorage) Flush() {
	if len(c.inputBag.old) > 0 && c.inputBag.isModified {
		oldInputStr, err := valuesToBase64(c.inputBag.old)

		if err != nil {
			iLogError(err.Error())
		} else {
			c.sessionData.Set("_old_input", oldInputStr)
		}
	}

	if len(c.inputBag.data) > 0 && c.inputBag.isModified {
		inputStr, err := valuesToBase64(c.inputBag.data)

		if err != nil {
			iLogError(err.Error())
		} else {
			c.sessionData.Set("_input", inputStr)
		}
	}

	if len(c.errorBag.err) > 0 && c.errorBag.isModified {
		errorStr, err := errorBagToBase64(c.errorBag)
		if err != nil {
			iLogError(err.Error())
		} else {
			c.sessionData.Set("_error_input", errorStr)
		}
	}
}

func (c *ContextSessionFlashStorage) Errors() *ErrorBag {
	return c.errorBag
}

func (c *ContextSessionFlashStorage) Inputs() *InputBag {
	return c.inputBag
}
