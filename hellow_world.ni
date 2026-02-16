var n = 50
var a = 0
var b = 1
var count = 0

while count < n{
    var temp = a
    a = b
    b = temp + b
    count = count + 1
}

print "50th fibonacci number:"
print a