package user_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Docs :- https://github.com/mongodb/mongo-go-driver

// go mod init <module-name>
// go mod tidy  //  Install all the packages dependency related to our code

// If both models are in the same package (e.g., models), you can directly reference Product_Model from user_model.

type User struct {
	ID              primitive.ObjectID `json:_id bson:_id`
	First_Name      *string
	Last_Name       *string
	Password        *string
	Email           *string
	Phone           *string
	Token           *string
	Refresh_Token   *string
	Created_At      time.Time
	Updated_At      time.Time
	User_ID         *string
	User_Cart       []ProductUser
	Address_Details []Address
	Order_Status    []Order
}

// ---- Reason to Use *string (Pointer String)

// Using `*string` (a pointer to a string) instead of `string` in the struct fields serves specific purposes related to flexibility, optionality, and memory efficiency. Below are the reasons and scenarios where `*string` is preferred over `string`:

// ---

// ### **1. Optional Fields (Nullable Values)**
// - In many cases, the fields in the struct (e.g., `First_Name`, `Email`, etc.) may not always have a value.
// - A Go `string` type cannot represent a `null` or `nil` value — it defaults to an empty string `""` if uninitialized.
// - By using `*string`, you can explicitly set a field to `nil`, making it clear that the value is absent or unknown.

// Example:
// ```go
// user := User{
//     First_Name: nil, // Indicates that the first name is not set
// }
// ```

// ---

// ### **2. Database Integration**
// - Many databases, such as MongoDB, support `null` values for fields. When working with such databases, using pointers (`*string`) allows the struct to align with the database schema.
// - For example:
//   - If a user's `Last_Name` is absent in MongoDB, it might be stored as `null`.
//   - To correctly map this in Go, `Last_Name` must be a pointer (`*string`), allowing the value to be `nil` (instead of an empty string).

// ---

// ### **3. Reduce Memory Usage for Large Data**
// - Passing around a `*string` uses less memory compared to a regular `string` because the pointer itself is smaller (typically 8 bytes) than a string that includes both the data and its metadata (e.g., length).
// - This is beneficial when you frequently copy large structs or manipulate large strings.

// ---

// ### **4. Differentiating Between "Empty" and "Unset"**
// - With a plain `string`, you can't differentiate between a field being intentionally set to `""` (empty string) and a field that has never been set (unset).
// - Using `*string` makes this distinction:
//   - `nil` → Field is unset or unknown.
//   - `""` → Field is explicitly set to an empty value.

// ---

// ### **5. Compatibility with JSON/BSON Marshalling**
// - When marshalling/unmarshalling structs to/from JSON or BSON (used by MongoDB), using `*string` ensures that `null` values are correctly handled.
// - For example, with `*string`:
//   ```json
//   {
//     "First_Name": null
//   }
//   ```
//   With `string`, `null` would often be converted to an empty string `""`, potentially misrepresenting the data.

// ---

// ### **6. Flexibility for Validation**
// - Using pointers makes it easier to validate if a field is set.
// - Example:
//   ```go
//   if user.First_Name == nil {
//       fmt.Println("First name is not provided")
//   }
//   ```

// ---

// ### **Considerations When Using `*string`**

// 1. **Dereferencing:**
//    Accessing the value requires dereferencing:
//    ```go
//    if user.First_Name != nil {
//        fmt.Println(*user.First_Name) // Dereference to get the actual string value
//    }
//    ```

// 2. **Increased Complexity:**
//    Working with pointers introduces additional complexity, such as handling `nil` values and potential bugs from improper dereferencing.

// 3. **Performance Trade-off:**
//    For small strings or non-nullable fields, using `string` directly may be simpler and slightly more performant.

// ---

// ### **When to Use `string` Instead of `*string`**

// Use plain `string` if:
// - The field is mandatory and will always have a value.
// - The concept of `null` or `nil` doesn't apply to the field in your context.
// - The extra memory and complexity of using pointers isn't justified.

// ---

// ### **Conclusion**
// In your `User` struct, fields like `First_Name`, `Last_Name`, `Email`, etc., are using `*string` because:
// - These fields might be optional or nullable in your application or database.
// - This allows you to differentiate between "unset" (`nil`) and "empty" (`""`) values.
// - It's particularly useful for database integrations, where `null` values are common.

// If these fields are guaranteed to always have values, you could switch to using `string` to simplify the code. However, the current design is more flexible and aligns well with MongoDB's BSON schema.
