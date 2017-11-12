package datamapper

import (
	"fmt"
	"github.com/juju/errors"
)

// Entity defines a domain entity
type Entity interface {
	// Return an entity's ID
	GetID() string
}

// UnitOfWork defines the persistence work needed to be accomplished for a
// given business transaction
type UnitOfWork struct {
	// A map of objects to be inserted into the datastore
	newObjects map[string]Entity

	// A map of objects to be updated in the datastore
	dirtyObjects map[string]Entity

	// A map of objects to be deleted from the datastore
	deletedObjects map[string]Entity
}

// NewUnitOfWork creates a new instance of UnitOfWork
func NewUnitOfWork() *UnitOfWork {
	return &UnitOfWork{
		newObjects:     make(map[string]Entity),
		dirtyObjects:   make(map[string]Entity),
		deletedObjects: make(map[string]Entity),
	}
}

// assertEntityHasAnID returns an error if an entity has no ID
func (unit *UnitOfWork) assertEntityHasID(entity Entity) error {
	if entity.GetID() == "" {
		return errors.Errorf(
			"Registering entity failed: entity has no ID: %+v",
			entity,
		)
	}

	return nil
}

// assertEntityNotRegisteredAs returns an error if an entity is already registered
// as having some state of persistence (or lack thereof)
func (unit *UnitOfWork) assertEntityNotRegisteredAs(entity Entity, registeredAs string) error {
	var registry map[string]Entity

	switch registeredAs {
	case "dirty":
		registry = unit.dirtyObjects
	case "deleted":
		registry = unit.deletedObjects
	case "new":
		registry = unit.newObjects
	default:
		panic(
			fmt.Sprintf(
				"unknown registry for state of persistence: \"%s\"",
				registeredAs,
			),
		)
	}

	_, entityAlreadyRegistered := registry[entity.GetID()]
	if entityAlreadyRegistered {
		return errors.Errorf(
			"Registering entity failed: entity with ID \"%s\" is already registered as %s",
			entity.GetID(),
			registeredAs,
		)
	}

	return nil
}

// RegisterNew registers a domain entity as being new
func (unit *UnitOfWork) RegisterNew(entity Entity) error {
	// Entity must have an ID
	if err := unit.assertEntityHasID(entity); err != nil {
		return err
	}

	// New entity must not be dirty
	if err := unit.assertEntityNotRegisteredAs(entity, "dirty"); err != nil {
		return err
	}

	// New entity must not be slated for deletion
	if err := unit.assertEntityNotRegisteredAs(entity, "deleted"); err != nil {
		return err
	}

	// New entity must not already be registered as new
	if err := unit.assertEntityNotRegisteredAs(entity, "new"); err != nil {
		return err
	}

	unit.newObjects[entity.GetID()] = entity

	return nil
}

// RegisterDirty registers a domain entity as being dirty
func (unit *UnitOfWork) RegisterDirty(entity Entity) error {
	// Entity must have an ID
	if err := unit.assertEntityHasID(entity); err != nil {
		return err
	}

	// Dirty entity must not be slated for deletion
	if err := unit.assertEntityNotRegisteredAs(entity, "deleted"); err != nil {
		return err
	}

	_, entityIsAlreadyDirty := unit.newObjects[entity.GetID()]
	if !entityIsAlreadyDirty {
		unit.dirtyObjects[entity.GetID()] = entity
	}

	return nil
}

// RegisterDeleted registers a domain entity as being deleted
func (unit *UnitOfWork) RegisterDeleted(entity Entity) error {
	// Entity must have an ID
	if err := unit.assertEntityHasID(entity); err != nil {
		return err
	}

	// No need to mark something for deletion if it is still new and hasn't
	// been persisted
	_, entityIsAlreadyNew := unit.newObjects[entity.GetID()]
	if entityIsAlreadyNew {
		delete(unit.newObjects, entity.GetID())

		return nil
	}

	// Remove an entity from the dirty set if it is to be deleted
	delete(unit.dirtyObjects, entity.GetID())

	_, entityIsAlreadyDeleted := unit.deletedObjects[entity.GetID()]
	if !entityIsAlreadyDeleted {
		unit.deletedObjects[entity.GetID()] = entity
	}

	return nil
}
