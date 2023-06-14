package ai.lepo.java;

import java.util.ArrayList;
import java.util.List;

public class SimpleClass {

    List<String> getUsers() {
        List<String> users = new ArrayList<>();

        users.add("Person1");
        users.add("Person2");

        return users;
    }


    void printMessage(final String message) {
        System.out.println(message);
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
     */
    void blockCommentedMethod(final Long id, final String message) {
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

    }


}
