package sorting

import (
	"reflect"
	"testing"
)

type MockIsGreaterer int

func (i MockIsGreaterer) IsGreater(e IsGreaterer) bool {
	return i > e.(MockIsGreaterer)
}

func streamEqual(stream Stream, expected ...int) bool {
	i := 0
	for v := range stream {
		if i >= len(expected) {
			return false
		}
		if MockIsGreaterer(expected[i]) != v {
			return false
		}
		i++
	}
	return i == len(expected)

}

func CreateStream(values ...int) Stream {
	s := make(Stream)
	go func() {
		for _, val := range values {
			s <- MockIsGreaterer(val)
		}
		close(s)
	}()

	return s
}

func TestReadFirstElementOfAllTheStreams(t *testing.T) {
	s1 := CreateStream(0, 1)
	arr := readFirstElementOfAllTheStreams([]Stream{s1})
	if !reflect.DeepEqual(arr, []IsGreaterer{MockIsGreaterer(0)}) {
		t.Error(arr)
	}

	if next := <-s1; next != MockIsGreaterer(1) {
		t.Error("read too many element in the list")
		return
	}

	s1 = CreateStream(0, 1, 2)
	s2 := CreateStream(5, 6, 7)
	s3 := CreateStream(90, 91, 92)
	arr = readFirstElementOfAllTheStreams([]Stream{s1, s2, s3})
	if !reflect.DeepEqual(arr, []IsGreaterer{MockIsGreaterer(0), MockIsGreaterer(5), MockIsGreaterer(90)}) {
		t.Error(arr)
	}
}

func TestAllStreamsButOne(t *testing.T) {
	streams := []Stream{CreateStream(0, 10, 1), CreateStream(1, 10, 1), CreateStream(2, 10, 1)}

	res := allStreamsButOne(streams, 1)
	if len(res) != 2 {
		t.Error()
	}
	if e := <-res[0]; e != MockIsGreaterer(0) {
		t.Error()
	}
	if e := <-res[1]; e != MockIsGreaterer(2) {
		t.Error()
	}
}

func TestGetSmallestElement(t *testing.T) {
	e, i := getSmallestElement([]IsGreaterer{MockIsGreaterer(9), MockIsGreaterer(10), MockIsGreaterer(1)})
	if i != 2 {
		t.Error("bad index")
	}
	if e != MockIsGreaterer(1) {
		t.Error("bad value")
	}
}

func TestGetSecondSmallestElement(t *testing.T) {
	e, i := getSecondSmallestElement([]IsGreaterer{MockIsGreaterer(9), MockIsGreaterer(10), MockIsGreaterer(1)})
	if i != 0 {
		t.Error("bad index")
	}
	if e != MockIsGreaterer(9) {
		t.Error("bad value")
	}

	e, i = getSecondSmallestElement([]IsGreaterer{MockIsGreaterer(0), MockIsGreaterer(10), MockIsGreaterer(1)})
	if i != 2 {
		t.Error("bad index")
	}
	if e != MockIsGreaterer(1) {
		t.Error("bad value")
	}
}

func TestGetSmallerStreamWithSecondSmallestElementWithAllTheOthersGreaterStreams(t *testing.T) {
	streams := []Stream{CreateStream(0, 1), CreateStream(3, 4), CreateStream(7, 8)}
	stream, e, streams := getSmallerStreamWithSecondSmallestElementWithAllTheOthersGreaterStreams(streams)
	var _ = stream
	_ = streams
	if e != MockIsGreaterer(3) {
		t.Error("getSmallerStreamWithSecondSmallestElementWithAllTheOthersGreaterStreams does not return the second smallest value", e)
	}

	if v := <-stream; v != MockIsGreaterer(0) {
		t.Error("return streams is not the smalest", v)
	}

	if len(streams) != 2 {
		t.Error("bad number of returned streams")
	}

	if len(streams) != 2 {
		t.Error("bad number of returned streams")
	}
}

func TestSortStreams(t *testing.T) {
	if s := SortStreams(CreateStream(0, 1, 2)); !streamEqual(s, 0, 1, 2) {
		t.Error()
	}

	if s := SortStreams(CreateStream(0, 1), CreateStream(2, 3)); !streamEqual(s, 0, 1, 2, 3) {
		t.Error()
	}

	if s := SortStreams(CreateStream(0), CreateStream(2)); !streamEqual(s, 0, 2) {
		t.Error()
	}

	if s := SortStreams(CreateStream(0, 5), CreateStream(2)); !streamEqual(s, 0, 2, 5) {
		t.Error()
	}

	if s := SortStreams(CreateStream(0, 5), CreateStream(2, 7)); !streamEqual(s, 0, 2, 5, 7) {
		t.Error()
	}

	if s := SortStreams(CreateStream(0, 5), CreateStream(2, 7, 9, 10, 12)); !streamEqual(s, 0, 2, 5, 7, 9, 10, 12) {
		t.Error()
	}

	if s := SortStreams(CreateStream()); !streamEqual(s) {
		t.Error()
	}

	if s := SortStreams(CreateStream(1, 4), CreateStream(0, 5)); !streamEqual(s, 0, 1, 4, 5) {
		t.Error()
	}

}
