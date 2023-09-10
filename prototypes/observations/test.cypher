CREATE (td:TypeDecl { name: "DishItem", qname: "com.example.app.models.DishItem", code: "data class DishItem(val name: String, val price: Int)
", parentqname: "", kind: "data_class"});

CREATE (td:TypeDecl { name: "Order", qname: "com.example.app.models.Order", code: "var orderId: Int,
val userId: Int,
val shopId: Int,
val items: List<DishItem>,
val amount: BigDecimal,
var status: OrderConfirmation
", parentqname: "", kind: "data_class"});

create (td:TypeDecl{name: "Shop", qname: "com.example.app.models.Shop", code: "data class Shop(
val shopId: Int,
val name: String,
var open: Boolean = false ,
val items: MutableList<DishItem> = mutableListOf(),
val orders: MutableList<Order> = mutableListOf()
)", parentqname: "", kind: "class"});

create (td:TypeDecl{name:"User", qname: "com.example.app.models.User", code:"data class User(
val id: Int,
val name: String,
val age: Int,
val orderHistory: MutableList<Order>
)", parentqname: "", kind: "class"});

create (td:TypeDecl{name: "Vendor", qname: "com.example.app.models.Vendor", code: "data class Vendor(val name: String, val shopId: Int)
", parentqname: "", kind: "class"} );

CREATE (
td:TypeDecl { name: "OrderConfirmation", qname: "com.example.app.models.OrderConfirmation", code: "enum class OrderConfirmation {
  CONFIRMED, REJECTED
  }", parentqname: "", kind: "class"
  });
  
  create (:TypeDecl{name: 'UserService', qname:'com.example.app.shop.UserService', code: "package com.example.app.user
  
  import com.example.app.models.Order
  import com.example.app.models.User
  import com.example.app.models.Vendor
  
  interface UserService {
    
    fun getUser(userId: Int): User?
    
    fun login(user: User): Boolean;
    
    fun addToUserHistory(userId: Int, order: Order): Boolean
    
    fun makeUserVendor(userId: Int): Vendor?
    }", parentqname: '', kind:'interface'})<-[:IMPLEMENTS]-(:TypeDecl{name: 'UserServiceImpl', qname: 'com.example.app.shop.UserServiceImpl', code: 'package com.example.app.user
    
    import com.example.app.models.Order
    import com.example.app.models.User
    import com.example.app.models.Vendor
    
    class UserServiceImpl(
    private val users: MutableList<User> = mutableListOf(
    User(id = 1, name = "Nithin", age = 23, mutableListOf()),
    User(id = 2, name = "Raju", age = 24, mutableListOf()),
    ),
    private val vendors: MutableList<Vendor> = mutableListOf()
    ) : UserService {
      override fun getUser(userId: Int): User? {
        RETURN users.firstOrNull { user -> user.id == userId }
      }
      
      override fun login(user: User): Boolean {
        RETURN user.name == "Nithin"
      }
      
      override fun addToUserHistory(userId: Int, order: Order): Boolean {
        val user = users.firstOrNull { user -> user.id == userId } ?:
        RETURN false
        user.orderHistory.add(order);
        
        RETURN true
      }
      
      override fun makeUserVendor(userId: Int): Vendor? {
        val user = users.firstOrNull { user -> user.id == userId } ?:
        RETURN null
        val vendor = Vendor(name = user.name, shopId = 1)
        vendors.add(vendor)
        RETURN vendor
      }
      
      }', kind:"class", parentqname:"com.example.app.shop.UserService"})
      
      create (m:Method{code: "override fun getUser(userId: Int): User? {
        RETURN users.firstOrNull { user -> user.id == userId }
        }", name: 'getUser', qname: 'com.example.app.shop.UserServiceImpl.getUser', parentqname: 'com.example.app.shop.UserServiceImpl'})
