package output

import (
	"fmt"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func TestCursor(t *testing.T) {
	list1 := &ListNode{Val: 0, Next: &ListNode{Val: 3}}
	list2 := &ListNode{Val: 0, Next: &ListNode{Val: 1, Next: &ListNode{Val: 2}}}

	traverseList(mergeTwoLists(list1, list2))
}

func traverseList(head *ListNode) {
	cursor := head
	for {
		if cursor == nil {
			break
		}

		fmt.Printf("%+v\n", cursor.Val)
		cursor = cursor.Next
	}
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil && list2 == nil {
		return nil
	}

	if list1 == nil {
		return list2
	}

	if list2 == nil {
		return list1
	}

	leftCursor := list1
	rightCursor := list2

	head := &ListNode{}
	cursor := head

	for {
		if leftCursor.Val < rightCursor.Val {
			cursor.Next = leftCursor
			cursor = cursor.Next
			if leftCursor.Next != nil {
				leftCursor = leftCursor.Next
			} else {
				cursor.Next = rightCursor
				break
			}
		} else {
			cursor.Next = rightCursor
			cursor = cursor.Next
			if rightCursor.Next != nil {
				rightCursor = rightCursor.Next
			} else {
				cursor.Next = leftCursor
				break
			}
		}
	}

	return head.Next
}
