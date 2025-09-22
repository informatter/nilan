var a = "hellow world!"
print a


var n = 20
var a = 0
var b = 1
var count = 0

while count < n{
    var temp = a
    a = b
    b = temp + b
    count = count + 1
}

print "20th fibonacci number:"
print a