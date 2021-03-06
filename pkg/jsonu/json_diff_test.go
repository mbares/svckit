package jsonu

import (
	"encoding/json"
	"github.com/minus5/svckit/pkg/testu"
	"testing"

	"github.com/minus5/go-simplejson"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	before string
	after  string
	diff   string
}

func testCases() []testCase {
	return []testCase{
		{
			before: `{"attr1":"v1","attr2":2,"attr3":3.3,"attr4":true}`,
			after:  `{"attr1":"v2","attr2":2,"attr3":3.3,"attr4":true}`,
			diff:   `{"attr1":"v2"}`,
		},
		{
			before: `{"attr1":"v1","attr2":2,"attr3":3.3,"attr4":true}`,
			after:  `{"attr1":"v1","attr2":3,"attr3":3.3,"attr4":true}`,
			diff:   `{"attr2":3}`,
		},
		{
			before: `{"attr1":"v1","attr2":2,"attr3":3.3,"attr4":true}`,
			after:  `{"attr1":"v1","attr2":3,"attr3":3.4,"attr4":true}`,
			diff:   `{"attr2":3,"attr3":3.4}`,
		},
		{
			before: `{"attr1":"v1","attr2":2,"attr3":3.3,"attr4":true}`,
			after:  `{"attr1":"v1","attr2":3,"attr3":3.4,"attr4":false}`,
			diff:   `{"attr2":3,"attr3":3.4,"attr4":false}`,
		},
		{ //atribut koji vise ne postoji postavlja se na null
			before: `{"attr1":"v1","attr3":3.3}`,
			after:  `{"attr1":"v1","attr2":3}`,
			diff:   `{"attr2":3,"attr3":null}`,
		},
		{ //array
			before: `{"arr1":[1,2,3]}`,
			after:  `{"arr1":[1,2,4]}`,
			diff:   `{"arr1":[1,2,4]}`,
		},
		{ //dva ista array-a
			before: `{"arr1":[1,2,3]}`,
			after:  `{"arr1":[1,2,3]}`,
			diff:   `{}`,
		},
		{ //object
			before: `{"o1":{"k1":"v1","k2":"v2"}}`,
			after:  `{"o1":{"k1":"v1","k2":"v3"}}`,
			diff:   `{"o1":{"k2":"v3"}}`,
		},
		{ //dva ista object-a
			before: `{"o1":{"k1":"v1","k2":"v2"}}`,
			after:  `{"o1":{"k1":"v1","k2":"v2"}}`,
			diff:   `{}`,
		},
		{ //object na drugom nivou
			before: `{"o1":{"k1":"v1","k2":{"k3":3,"k4":5}}}`,
			after:  `{"o1":{"k1":"v1","k2":{"k3":3,"k4":6,"k5":5}}}`,
			diff:   `{"o1":{"k2":{"k4":6,"k5":5}}}`,
		},
		{ //object na drugom nivou
			before: `{"o1":{"k1":"v1","k2":{"-1":3,"-2":5}}}`,
			after:  `{"o1":{"k1":"v1","k2":{"-1":3,"-2":6,"-3":5}}}`,
			diff:   `{"o1":{"k2":{"-2":6,"-3":5}}}`,
		},
		{ //dva object key-a
			before: `{"o1":{"k1":"v1","k2":"v2"},"o2":{"k3":"v3"}}`,
			after:  `{"o1":{"k1":"v1","k2":"v3"},"o2":{"k3":"v4"}}`,
			diff:   `{"o1":{"k2":"v3"},"o2":{"k3":"v4"}}`,
		},
	}
}

func TestJsonDiff(t *testing.T) {
	for _, d := range testCases() {
		//log.Printf("%v", d)
		diff := diff([]byte(d.before), []byte(d.after))
		assert.Equal(t, string(diff), d.diff)

		full := merge([]byte(d.before), []byte(d.diff))
		assert.Equal(t, string(full), d.after)
	}
}

func TestMerge0(t *testing.T) {
	//d := testCases()[0]
	full := merge(
		[]byte(`{"o1":{"k1":"v1","k2":"v2"},"o3":{}}`),
		[]byte(`{"o1":{"k2":"v3"},"o2":{"k3":4},"o3":{"k4":5}}`))
	assert.Equal(t, `{"o1":{"k1":"v1","k2":"v3"},"o2":{"k3":4},"o3":{"k4":5}}`, string(full))
}

func TestMerge00(t *testing.T) {
	full := merge(
		[]byte(`{"k":"v"}`),
		[]byte(`{"k":"v2"}`))
	assert.Equal(t, `{"k":"v2"}`, string(full))
}

func TestMergeCopyMap(t *testing.T) {
	fe := `{"k1":"v1","o1":{"k2":"v2"},"o2":{"k3":1,"o3":{"o4":"v3"}}}`
	f, err := simplejson.NewJson([]byte(fe))
	assert.Nil(t, err)
	d, err := simplejson.NewJson([]byte(`{"k1":"v11","o1":{"k2":"v22"}}`))
	assert.Nil(t, err)
	f2 := Merge(f, d)
	buf, _ := f2.Encode()
	e := `{"k1":"v11","o1":{"k2":"v22"},"o2":{"k3":1,"o3":{"o4":"v3"}}}`
	assert.Equal(t, e, string(buf))

	// f nije promjenjen nakon merge
	buf, _ = f.Encode()
	assert.Equal(t, fe, string(buf))

	//fo2 := f.MustMap()["o2"].(map[string]interface{})["o3"]
	//f2o2 := f2.MustMap()["o2"].(map[string]interface{})["o3"]
	//fmt.Printf("%v %v\n", &fo2, &f2o2)
	//assert.True(t, &fo2 == &f2o2)
	//assert.Equal(t, f.Get("o1"), f2.Get("o1"))

}

func TestMerge1(t *testing.T) {
	f, err := simplejson.NewJson([]byte(`{}`))
	assert.Nil(t, err)
	d, err := simplejson.NewJson([]byte(`{"o":1}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ := f.Encode()
	assert.Equal(t, `{"o":1}`, string(buf))

	d, err = simplejson.NewJson([]byte(`{"o2":2}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ = f.Encode()
	assert.Equal(t, `{"o":1,"o2":2}`, string(buf))

	d, err = simplejson.NewJson([]byte(`{"o3":{"k1":5}}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ = f.Encode()
	assert.Equal(t, `{"o":1,"o2":2,"o3":{"k1":5}}`, string(buf))

	d, err = simplejson.NewJson([]byte(`{"o3":{"k1":1}}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ = f.Encode()
	assert.Equal(t, `{"o":1,"o2":2,"o3":{"k1":1}}`, string(buf))

	d, err = simplejson.NewJson([]byte(`{"o3":{"k1":2,"o4":{"k2":1}}}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ = f.Encode()
	assert.Equal(t, `{"o":1,"o2":2,"o3":{"k1":2,"o4":{"k2":1}}}`, string(buf))

	d, err = simplejson.NewJson([]byte(`{"o3":{"o4":{"k2":2}}}`))
	assert.Nil(t, err)
	f = Merge(f, d)
	buf, _ = f.Encode()
	assert.Equal(t, `{"o":1,"o2":2,"o3":{"k1":2,"o4":{"k2":2}}}`, string(buf))

}

func TestSameKeyIntArray(t *testing.T) {
	m1 := map[string]interface{}{"k": []int{1, 2, 3}}
	m2 := map[string]interface{}{"k": []int{1, 2, 4}}
	m3 := map[string]interface{}{"k": []int{1, 2, 3}}
	j1 := MapToSimplejson(m1)
	j2 := MapToSimplejson(m2)
	j3 := MapToSimplejson(m3)
	assert.True(t, sameKey("k", j1, j2) == areDifferent)
	assert.True(t, sameKey("k", j1, j3) == areSame)
}

func TestJsonMerge(t *testing.T) {
	for _, c := range testCases() {
		var before, diff map[string]interface{}
		json.Unmarshal([]byte(c.before), &before)
		json.Unmarshal([]byte(c.diff), &diff)
		JsonMerge(before, diff)
		aa, _ := json.Marshal(before)
		afterActual := string(aa)
		if afterActual != c.after {
			t.Logf("actual:   %s", afterActual)
			t.Logf("expected: %s", c.after)
		}
		assert.Equal(t, afterActual, c.after)
	}
}

func TestDeepCopyMap(t *testing.T) {
	s := `{"o1":{"k1":"v1","k2":"v2"},"o2":{"k3":"v3"}}`
	var m map[string]interface{}
	json.Unmarshal([]byte(s), &m)
	m2 := DeepCopyMap(m)
	s2, _ := json.Marshal(m2)
	assert.Equal(t, s, string(s2))
}

func TestMapEqual(t *testing.T) {
	inner := &map[string]interface{}{"k": "v"}
	m1 := map[string]interface{}{"k": inner}
	m2 := map[string]interface{}{"k": inner}
	//m2 := m1
	k1 := m1["k"]
	k2 := m1["k"]
	//fmt.Println(k1, k2)
	assert.True(t, k1 == k2)
	j1 := MapToSimplejson(m1)
	j2 := MapToSimplejson(m2)
	j1m := j1.Get("k").Interface()
	j2m := j2.Get("k").Interface()
	//fmt.Println(j1m, j2m)
	assert.True(t, j1m == j2m)
}

func TestDiffObjectPointerNaIstiMap(t *testing.T) {
	inner := &map[string]interface{}{"k0": "v"}
	l := map[string]interface{}{"k1": inner}
	r := map[string]interface{}{"k1": inner, "k2": 1}
	d := diffmap(l, r)
	assert.Equal(t, 1, len(d))
}

func TestDiffObjectPointerNaRazlicitMap(t *testing.T) {
	inner := &map[string]interface{}{
		"k0": "v",
		"km": &map[string]interface{}{"k4": 4},
	}
	inner2 := &map[string]interface{}{
		"k0": "v",
		"km": &map[string]interface{}{"k4": 5},
	}
	l := map[string]interface{}{"k1": inner}
	r := map[string]interface{}{"k1": inner2, "k2": 1}
	d := diffmap(l, r)
	assert.Equal(t, 2, len(d))
	//fmt.Printf("%#v\n", d)
}

func TestBugFix(t *testing.T) {
	t.Skip("vizualno testiranje")
	full := `{
	"S": {
		" X": {
			"C": {
				" )%": {
					"T": {
						" g)": {
							"M": {
								"%Fh;": {
									"#1": {
										"en": "St Kilda Saints",
										"hr": "St Kilda Saints"
									},
									"#2": {
										"en": "Greater Western Sydney",
										"hr": "Greater Western Sydney"
									},
									"1": "St Kilda Saints",
									"2": "Greater Western Sydney",
									"O": {
										" !": {
											"1": 42,
											"2": 35,
											"y": "Current"
										}
									},
									"a": 20,
									"d": "2017-05-05T11:50:00+02:00",
									"i": 11279535,
									"u": "2017-05-05T11:50:37+02:00",
									"v": {
										"#c": {
											"en": "Melbourne",
											"hr": "Melbourne"
										},
										"#o": {
											"en": "Australia",
											"hr": "Australija"
										},
										"#t": {
											"en": "Etihad Stadium",
											"hr": "Etihad"
										},
										"c": "Melbourne",
										"o": "Australia",
										"t": "Etihad Stadium"
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
`
	diff := `{"S":{" X":{"C":{" )%":{"T":{" g)":{"M":{"%Fh;":{"O":{" !":{"1":43,"2":48}}}}}}}}}}}`
	diff2 := `{"S":{" X":{"C":{" )%":{"T":{" g)":{"M":{"%Fh;":{"O":{" !":{"1":44,"2":49}}}}}}}}}}}`
	f, err := simplejson.NewJson([]byte(full))
	assert.Nil(t, err)
	d, err := simplejson.NewJson([]byte(diff))
	assert.Nil(t, err)
	d2, err := simplejson.NewJson([]byte(diff2))
	assert.Nil(t, err)

	f2 := Merge(f, d)
	f2 = Merge(f, d2)
	//f2 := Diff(d, f)
	buf, err := f2.EncodePretty()
	assert.Nil(t, err)
	testu.PPBuf(buf)
}
