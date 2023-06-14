package ai.lepo.java;

import java.net.http.HttpClient;
import java.util.List;

public interface SimpleInterface {

    List<String> getNames();

    String greet(final String message);


    // This is a line comment
    void lineCommentMethod(String something);


    /**
     * This is a block comment that explains some things about this getHttpClient method
     *
     * @param address
     * @return HttpClient
     */
    HttpClient getHttpClient(final String address);

}
