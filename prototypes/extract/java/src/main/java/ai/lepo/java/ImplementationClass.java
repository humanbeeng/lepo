package ai.lepo.java;

import java.net.http.HttpClient;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class ImplementationClass implements SimpleInterface {
    @Override
    public List<String> getNames() {
        var someList = new ArrayList<String>();
        someList.add("Hello");
        someList.add("there");
        return someList;
    }

    @Override
    public String greet(String message) {
        System.out.println("I am greeting" + message);
        return message;
    }

    @Override
    public void lineCommentMethod(String something) {
        // do nothing
        var a = 1;
    }

    @Override
    public HttpClient getHttpClient(String address) {
        return null;
    }
}
