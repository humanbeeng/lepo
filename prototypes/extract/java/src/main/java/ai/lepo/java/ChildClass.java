package ai.lepo.java;

import java.util.HashSet;
import java.util.Set;

public class ChildClass extends ParentClass {

    @Override
    void overridableMethod(final String message) {
        // Haha, you just got overridden
        Set<String> someSet = new HashSet<>();
    }
}
