package datamapper

import (
	"strings"
	"testing"
)

type EntityStub struct {
	id string
}

func (m *EntityStub) GetID() string {
	return m.id
}

func errShouldContainStr(err error, s string, t *testing.T) {
	if !strings.Contains(err.Error(), s) {
		t.Errorf("error should contain string \"%s\", contains \"%s\"", s, err.Error())
	}
}

func failOnUnexpectedErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("unexpected error: \"%s\"", err.Error())
	}
}

func TestNewUnitOfWork(t *testing.T) {
	u := NewUnitOfWork()

	if len(u.newObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.newObjects)
	}

	if len(u.dirtyObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.dirtyObjects)
	}

	if len(u.deletedObjects) != 0 {
		t.Errorf("should have zero new objects registered\n%+v", u.deletedObjects)
	}
}

func TestAssertEntityNotRegisteredAsPanicsOnUknownState(t *testing.T) {
	u := NewUnitOfWork()

	m := &EntityStub{id: "5"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when unkown persistence state passed")
		}
	}()

	_ = u.assertEntityNotRegisteredAs(m, "some-state")
}

func TestRegisterNew(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	foundEntity, isRegisteredNew := u.newObjects["5"]

	if !isRegisteredNew {
		t.Error("should be registered as new")
	}

	if foundEntity != m {
		t.Errorf(
			"registered entity should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundEntity,
		)
	}
}

func TestRegisterNewFailsWithNoId(t *testing.T) {
	m := &EntityStub{}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the entity has no ID")
	}

	errShouldContainStr(err, "entity has no ID", t)
}

func TestRegisterNewFailsIfEntityDirty(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the entity is already registered as dirty")
	}

	errShouldContainStr(err, "already registered as dirty", t)
}

func TestRegisterNewFailsIfEntityDeleted(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDeleted(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the entity is already registered as deleted")
	}

	errShouldContainStr(err, "already registered as deleted", t)
}

func TestRegisterNewFailsIfEntityAlreadyNew(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterNew(m)

	if err == nil {
		t.Error("should return an error if the entity is already registered as new")
	}

	errShouldContainStr(err, "already registered as new", t)
}

func TestRegisterDirty(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	foundEntity, isRegisteredDirty := u.dirtyObjects["5"]

	if !isRegisteredDirty {
		t.Error("should be registered as dirty")
	}

	if foundEntity != m {
		t.Errorf(
			"registered entity should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundEntity,
		)
	}
}

func TestRegisterDirtyFailsWithNoId(t *testing.T) {
	m := &EntityStub{}

	u := NewUnitOfWork()

	err := u.RegisterDirty(m)

	if err == nil {
		t.Error("should return an error if the entity has no ID")
	}

	errShouldContainStr(err, "entity has no ID", t)
}

func TestRegisterDirtyDoesNothingIfEntityNew(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	foundEntity, isRegisteredNew := u.newObjects[m.GetID()]
	if !isRegisteredNew {
		t.Error("entity should be registered as new")
	}

	if foundEntity != m {
		t.Errorf(
			"registered entity should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundEntity,
		)
	}

	err = u.RegisterDirty(m)
	failOnUnexpectedErr(err, t)

	_, isRegisteredDirty := u.dirtyObjects[m.GetID()]
	if isRegisteredDirty {
		t.Error("entity should not be registered as dirty if already registered as new")
	}

	foundEntity, isStillRegisteredNew := u.newObjects[m.GetID()]
	if !isStillRegisteredNew {
		t.Error("entity should still be registered as new")
	}

	if foundEntity != m {
		t.Errorf(
			"registered entity should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundEntity,
		)
	}
}

func TestRegisterDirtyFailsIfEntityDeleted(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDeleted(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterDirty(m)

	if err == nil {
		t.Error("should return an error if the entity is already registered as deleted")
	}

	errShouldContainStr(err, "already registered as deleted", t)
}

func TestRegisterDeleted(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterDeleted(m)
	failOnUnexpectedErr(err, t)

	foundEntity, isRegisteredDeleted := u.deletedObjects["5"]

	if !isRegisteredDeleted {
		t.Error("should be registered as deleted")
	}

	if foundEntity != m {
		t.Errorf(
			"registered entity should be correct: \nexpected: %+v \nactual: %+v",
			m,
			foundEntity,
		)
	}
}

func TestRegisterDeletedFailsWithNoId(t *testing.T) {
	m := &EntityStub{}

	u := NewUnitOfWork()

	err := u.RegisterDeleted(m)

	if err == nil {
		t.Error("should return an error if the entity has no ID")
	}

	errShouldContainStr(err, "entity has no ID", t)
}

func TestRegisterDeletedDeletesNewObject(t *testing.T) {
	m := &EntityStub{id: "5"}

	u := NewUnitOfWork()

	err := u.RegisterNew(m)
	failOnUnexpectedErr(err, t)

	err = u.RegisterDeleted(m)

	_, isStillRegisteredNew := u.newObjects["5"]
	if isStillRegisteredNew {
		t.Error("should delete registration for a new object when registering that new object as deleted")
	}

	_, isRegisteredDeleted := u.deletedObjects["5"]
	if isRegisteredDeleted {
		t.Error("should not register an object that was previously new as slated for removal - new object has not been persisted so nothing to delete")
	}
}
