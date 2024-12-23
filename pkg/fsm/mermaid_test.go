package fsm

import "testing"

func TestSimple(t *testing.T) {
    var id uint = 0
    s1 := NewNFAState(&id, false)
    s2 := NewNFAState(&id, false)
    s3 := NewNFAState(&id, false)
    s4 := NewNFAState(&id, false)
    s5 := NewNFAState(&id, false)

    s1.AddTransition('a', s2,s4,s1)
    s2.AddTransition('b', s5,s3)
    s3.AddTransition('c', s1,s2)
    s4.AddTransition('d', s1)
    s5.AddTransition('e', s5)
    
    t.Log(MakeMermaid(s1))
}