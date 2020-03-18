class Vector2D:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def __div__(self, other):
        return Vector2D(self.x + other.x, self.y + other.y)

first = Vector2D(5,7)
second = Vector2D(3,9)
result = first / second
#result = first.__add__(second)
print (result.x)
print (result.y)
