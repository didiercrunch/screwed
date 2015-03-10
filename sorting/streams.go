package sorting

import "fmt"

var _ = fmt.Print

type IsGreaterer interface {
	// return true iff "this" is strictly greater than e
	IsGreater(e IsGreaterer) bool
}

type alwaysBigger struct{}

func (a *alwaysBigger) IsGreater(e IsGreaterer) bool {
	return e != nil
}

type Stream chan IsGreaterer

type streamOutput struct {
	Val IsGreaterer
	Idx int
}

func readFirstElementOfAllTheStreams(streams []Stream) []IsGreaterer {
	values := make([]IsGreaterer, len(streams))
	ch := make(chan *streamOutput)
	for i, _ := range values {
		go func(idx int) {
			v := <-streams[idx]
			ch <- &streamOutput{v, idx}

		}(i)
	}

	i := 0
	for streamOutput := range ch {
		values[streamOutput.Idx] = streamOutput.Val
		i++
		if i == len(streams) {
			return values
		}
	}
	return values
}

func allStreamsButOne(streams []Stream, idx int) []Stream {
	ret := make([]Stream, len(streams)-1)
	for i, s := range streams {
		if i < idx {
			ret[i] = s
		} else if i > idx {
			ret[i-1] = s
		}
	}
	return ret
}

func prependToStream(streams Stream, toPrepend IsGreaterer) Stream {
	ret := make(Stream)
	go func() {
		ret <- toPrepend
		for e := range streams {
			ret <- e
		}
		close(ret)
	}()
	return ret
}

func prependValueToEachStream(streams []Stream, values []IsGreaterer) []Stream {
	ret := make([]Stream, len(streams))
	for i := 0; i < len(streams); i++ {
		ret[i] = prependToStream(streams[i], values[i])

	}
	return ret
}

func getSmallerStreamWithSecondSmallestElementWithAllTheOthersGreaterStreams(streams []Stream) (Stream, IsGreaterer, []Stream) {
	firstElements := readFirstElementOfAllTheStreams(streams)
	streams = prependValueToEachStream(streams, firstElements)

	var smallestStreamIndex int = -1
	var smallestItem IsGreaterer = &alwaysBigger{}
	var secondSmallestItem IsGreaterer = &alwaysBigger{}

	for i, e := range firstElements {
		if smallestItem.IsGreater(e) {
			smallestItem = e
			smallestStreamIndex = i
		} else if secondSmallestItem.IsGreater(e) {
			secondSmallestItem = e
		}
	}
	if _, ok := secondSmallestItem.(*alwaysBigger); ok {
		return streams[smallestStreamIndex], nil, []Stream{}
	}
	return streams[smallestStreamIndex], secondSmallestItem, allStreamsButOne(streams, smallestStreamIndex)
}

func sortStreamInChannel(o Stream, streams []Stream) {
	if len(streams) == 1 {
		for e := range streams[0] {
			o <- e
		}
		close(o)
		return
	}
	smallestStream, secondSmallestItem, greaterStreams := getSmallerStreamWithSecondSmallestElementWithAllTheOthersGreaterStreams(streams)
	if secondSmallestItem == nil {
		sortStreamInChannel(o, append(greaterStreams, smallestStream))
		return

	}
	for elm := range smallestStream {
		if !elm.IsGreater(secondSmallestItem) {
			o <- elm
		} else {
			smallestStream = prependToStream(smallestStream, elm)
			sortStreamInChannel(o, append(greaterStreams, smallestStream))
			return
		}
	}
	sortStreamInChannel(o, greaterStreams)
}

//  Takes sorted streams as input and return the merged result stream
func SortStreams(streams ...Stream) Stream {
	ret := make(Stream)
	go sortStreamInChannel(ret, streams)
	return ret
}
