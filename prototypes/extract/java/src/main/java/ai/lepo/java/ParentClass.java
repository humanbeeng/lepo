package ai.lepo.java;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;


// This is some info about the class


public class ParentClass {
    // this is some info about the class


    private final String privateMember = "I am a private member";

    public static String staticString = "I am a static string";

    private String instantiableMember;

    private Long anotherInstantiableMember;


    public ParentClass(final String instantiableMember) {
        this.instantiableMember = instantiableMember;
    }

    public ParentClass() {
        // no args constructor
    }

    /**
     * Block comment about the constructor
     *
     * @param anotherInstantiableMember
     */
    public ParentClass(final Long anotherInstantiableMember) {
        this.anotherInstantiableMember = anotherInstantiableMember;
    }


    public List<String> getUsers() {
        List<String> users = new ArrayList<>();

        users.add("Person1");
        users.add("Person2");

        return users;
    }


    void printMessage(final String message) {
        System.out.println(message);
    }

    void overridableMethod(final String someString) {
        for (int i = 0; i < 10; i++) {
            System.out.println("I am doing something for the " + i + "th time");
        }
    }


    // This is a single line comment about the below method.
    void singleLineCommentedMethod() {
        // do something about it.
    }


    /**
     * Send message to the user based on ID.
     *
     * @param id
     * @param message
     * @return BigDecimal
     */
    BigDecimal blockCommentedMethod(final Long id, final String message) {
        // do something about it.
        switch (message) {
            case "hello": {
                System.out.println(message);
                break;
            }
            default: {
                System.out.println("Default message");
            }
        }
        return BigDecimal.valueOf(1);
    }

    private void privateMethod() {
        // I am a private method
    }

    protected void protectedMethod() {
        // I am a protected method
    }

}
