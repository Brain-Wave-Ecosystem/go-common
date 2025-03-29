package dbx

import (
	"sort"

	"go.uber.org/zap"
)

type (
	ChangeRelationFunc[S any]               func(S) error
	ChangePositionsFunc[Q Queryable, S any] func(Q, S) error
)

type Keyer interface {
	Key() any
}

type Positioner interface {
	GetPosition() uint32
	SetPosition(pos uint32)
}

func AdjustRelation[S interface {
	Keyer
	comparable
}](
	prev, next []S,
	addFn, removeFn ChangeRelationFunc[S],
) error {
	prevM, nextM := toMap(prev), toMap(next)
	for key, prevItem := range prevM {
		if nextItem, contains := nextM[key]; !contains || prevItem != nextItem {
			err := removeFn(prevItem)
			if err != nil {
				return err
			}
		}
	}

	for key, nextItem := range nextM {
		if prevItem, contains := prevM[key]; !contains || prevItem != nextItem {
			err := addFn(nextItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ReassignPositions[S ~[]E, E Positioner, Q Queryable](q Q, elements S, updateFn ChangePositionsFunc[Q, E]) error {
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].GetPosition() < elements[j].GetPosition()
	})

	logger := zap.NewExample()
	l := zap.NewStdLog(logger)

	firstBrokenIndex := -1
	for i := 1; i < len(elements); i++ {
		if elements[i].GetPosition() != elements[i-1].GetPosition()+1 {
			firstBrokenIndex = i
			break
		}
	}

	if firstBrokenIndex == -1 {
		return nil
	}

	for i := firstBrokenIndex; i < len(elements); i++ {
		elements[i].SetPosition(elements[i-1].GetPosition() + 1)
	}

	for _, element := range elements {
		l.Printf("position: %d", element.GetPosition())

		err := updateFn(q, element)
		if err != nil {
			return err
		}
	}

	return nil
}

func toMap[S Keyer](slice []S) map[any]S {
	m := make(map[any]S, len(slice))
	for _, item := range slice {
		m[item.Key()] = item
	}
	return m
}
