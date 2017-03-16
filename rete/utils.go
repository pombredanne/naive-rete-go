package rete

import (
	"container/list"
	"github.com/beevik/etree"
	"errors"
)

func contain(l *list.List, value interface{}) *list.Element {
	if l == nil {
		return nil
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return e
		}
	}
	return nil
}

func remove_by_value(l *list.List, value interface{}) bool {
	if e := contain(l, value); e != nil {
		l.Remove(e)
		return true
	}
	return false
}

func FromXML(s string) (result []Production, err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(s)
	if err != nil {
		return result, err
	}
	root := doc.Root()
	if root == nil {
		return result, errors.New("Not XML")
	}

	for _, ep := range root.ChildElements() {
		if ep.Tag != "production" {continue}
		p := Production{
			rhs: make(map[string]interface{}),
		}
		for idx, hand := range ep.ChildElements() {
			if idx == 0 {
				p.lhs = parse_lhs(hand)
			} else if idx == 1 {
				for _, attr := range hand.Attr {
					p.rhs[attr.Key] = attr.Value
				}
			}
		}
		result = append(result, p)
	}
	return result, nil
}

func parse_lhs(root *etree.Element) Rule {
	r := NewRule()
	for _, e := range root.ChildElements() {
		switch e.Tag {
		case "has":
			class_name, identity, attribute, value := "", "", "", ""
			for _, attr := range e.Attr {
				if attr.Key == "classname" {
					class_name = attr.Value
				} else if attr.Key == "identifier" {
					identity = attr.Value
				} else if attr.Key == "attribute" {
					attribute = attr.Value
				} else if attr.Key == "value" {
					value = attr.Value
				}
			}
			has := NewHas(class_name, identity, attribute, value)
			r.items = append(r.items, has)
		case "filter":
			f := Filter{tmpl: e.Text()}
			r.items = append(r.items, f)
		}
	}
	return r
}
