## 0001

```go
func twoSum(nums []int, target int) []int {
	hash := make(map[int]int, len(nums))
	for i,v := range nums {
		if j, ok := hash[target-v]; ok {
			return []int{j, i}
		}
		hash[v] = i
	}
	return nil
}
```

