package datamapper

import (
	"github.com/juju/errors"
)

// Model defines a domain model
type Model interface {
	// Return a model's ID
	GetID() string
}

// UnitOfWork defines the persistence work needed to be accomplished for a
// given business transaction
type UnitOfWork struct {
	// A map of objects to be inserted into the datastore
	newObjects map[string]Model

	// A map of objects to be updated in the datastore
	dirtyObjects map[string]Model

	// A map of objects to be deleted from the datastore
	removedObjects map[string]Model
}

// NewUnitOfWork creates a new instance of UnitOfWork
func NewUnitOfWork() *UnitOfWork {
	return &UnitOfWork{
		newObjects:     make(map[string]Model),
		dirtyObjects:   make(map[string]Model),
		removedObjects: make(map[string]Model),
	}
}

// RegisterNew registers a domain model as being new
func (unit *UnitOfWork) RegisterNew(model Model) error {
	if model.GetID() == "" {
		return errors.Errorf(
			"Registering new model failed: model has no ID: %+v",
			model,
		)
	}

	_, modelIsDirty := unit.dirtyObjects[model.GetID()]
	if modelIsDirty {
		return errors.Errorf(
			"Registering new model failed: model with ID \"%s\" is already registered as dirty",
			model.GetID(),
		)
	}

	_, modelIsRemoved := unit.removedObjects[model.GetID()]
	if modelIsRemoved {
		return errors.Errorf(
			"Registering new model failed: model with ID \"%s\" is already registered as removed",
			model.GetID(),
		)
	}

	_, modelIsAlreadyNew := unit.newObjects[model.GetID()]
	if modelIsAlreadyNew {
		return errors.Errorf(
			"Registering new model failed: model with ID \"%s\" is already registered as new",
			model.GetID(),
		)
	}

	unit.newObjects[model.GetID()] = model

	return nil
}

// RegisterDirty registers a domain model as being dirty
func (unit *UnitOfWork) RegisterDirty(model Model) error {
	if model.GetID() == "" {
		return errors.Errorf(
			"Registering new model failed: model has no ID: %+v",
			model,
		)
	}

	_, modelIsRemoved := unit.removedObjects[model.GetID()]
	if modelIsRemoved {
		return errors.Errorf(
			"Registering dirty model failed: model with ID \"%s\" is already registered as removed",
			model.GetID(),
		)
	}

	_, modelIsAlreadyDirty := unit.newObjects[model.GetID()]
	if !modelIsAlreadyDirty {
		unit.dirtyObjects[model.GetID()] = model
	}

	return nil
}

// RegisterRemoved registers a domain model as being removed
func (unit *UnitOfWork) RegisterRemoved(model Model) error {
	if model.GetID() == "" {
		return errors.Errorf(
			"Registering removed model failed: model has no ID: %+v",
			model,
		)
	}

	// No need to mark something for deletion if it is still new and hasn't
	// been persisted
	_, modelIsAlreadyNew := unit.newObjects[model.GetID()]
	if modelIsAlreadyNew {
		delete(unit.newObjects, model.GetID())

		return nil
	}

	// Remove a model from the dirty set if it is to be deleted
	delete(unit.dirtyObjects, model.GetID())

	_, modelIsAlreadyRemoved := unit.removedObjects[model.GetID()]
	if !modelIsAlreadyRemoved {
		unit.removedObjects[model.GetID()] = model
	}

	return nil
}
