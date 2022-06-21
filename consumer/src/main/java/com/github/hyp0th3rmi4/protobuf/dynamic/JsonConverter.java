package com.github.hyp0th3rmi4.protobuf.dynamic;

import java.io.IOException;
import java.io.InputStream;
import java.io.FileWriter;
import java.io.ByteArrayOutputStream;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.google.protobuf.DescriptorProtos;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.cloudevents.jackson.JsonFormat;
import io.cloudevents.jackson.JsonCloudEventData;


/**
 * JsonConverter converts the payload of the cloud event
 * that is encoded in Protobuf into a corresponding JSON
 * structure by relying upon the information stored in the
 * associated file descriptor.
 */
public class JsonConverter {

    /**
     * A {@link JsonFormat} instance used to serialise/deserialise
     * the {@link CloudEvent} instance to and from JSON format.
     */
    protected JsonFormat format = new JsonFormat();

    /**
     * Converts the payload of the Cloud Event structure from Protobuf 
     * (binary) to JSON. The implementation reads the supplied file and
     * deserialises it into a Cloud Event instance in JSON format. Then,
     * based on the information stored in the "dataschema" attribute it
     * retrieves the associated file descriptor and uses this information
     * to dynamically parse the content of the payload and convert it 
     * into a corresponding JSON structure. If a target path is specified
     * it saves the rendered JSON structure to file, otherwise it dumps 
     * it to the standard output.
     * 
     * @param sourcePath    a {@link String} storing the full path to the
     *                      file containing the Cloud Event instance.
     * 
     * @param targetPath    a {@link String} storing the full path to the
     *                      file where the rendered JSON will be saved. It
     *                      can be {@literal null} in such case the rendered
     *                      JSON is written to the standard output. If the
     *                      file exists, it is rewritten.
     */
    public void convert(String sourcePath, String targetPath) throws Exception {

        // get the CloudEvebt and the descriptor set
        CloudEvent cloudEvent = this.getCloudEvent(sourcePath);
        DescriptorProtos.FileDescriptorSet descriptorSet = this.resolveDescriptorSet(cloudEvent.getDataSchema());

        // unpack the binary protobuf and convert it
        // to a corresponding payload expressed in 
        // JSON.
        String messageType = cloudEvent.getDataSchema().getFragment();
        byte[] payload = cloudEvent.getData().toBytes();
        JsonNode jsonPayload = this.convertPayload(descriptorSet, messageType, payload);
        JsonCloudEventData newPayload = JsonCloudEventData.wrap(jsonPayload);

        // regenerate the Cloud Event with the new payload.
        
        CloudEvent jsonPayloadCloudEvent = this.repackageEvent(cloudEvent, newPayload);

        ObjectWriter jsonWriter = new ObjectMapper().writerWithDefaultPrettyPrinter();
        String eventString = jsonWriter.writeValueAsString(jsonPayloadCloudEvent);

        // write the CloudEvent to the selected output.
        if (targetPath == null) {
            try (FileWriter writer = new FileWriter(targetPath, true)) {

                writer.write(eventString);
            }
        } else {    
            System.out.println(eventString);
        }

    }

    /**
     * Reads the content of the file identified by <i>sourcePath</i> and 
     * deserialises it into a {@link CloudEvent} instance and returns it
     * to the caller.
     * 
     * @param sourcePath    a {@link String} containing the path to the
     *                      file containing the JSON representation of
     *                      the cloud event.
     * 
     * @return  a {@link CloudEvent} instance representing the deserialised
     *          in-memory representation of the cloud event stored in the
     *          file passed as argument.
     */
    protected CloudEvent getCloudEvent(String sourcePath) throws IOException {

        Path path = Paths.get(sourcePath);
        byte[] buffer = Files.readAllBytes(path);

        JsonFormat format = new JsonFormat();
        return format.deserialize(buffer);
    }

    /**
     * Resolves the instance of {@link DescriptorProtos.FileDescriptorSet} from the 
     * data schema URI that is passed as argument to the method. The implementation
     * opens a connection to the supplied URI and reads all its content into an array
     * of bytes. It then uses {@link DescriptorProtos.FileDescriptorSet#parseFrom(byte[])}
     * to return an instace of the descriptor set.
     * 
     * @param dataSchemaUri a {@link URI} instance that wraps a reference to a resource
     *                      defining the file descriptor set.
     * 
     * @return  a {@link DescriptorProtos.FileDescriptorSet} instance representing the
     *          file descriptor pointed by the given URI.
     */
    protected DescriptorProtos.FileDescriptorSet resolveDescriptorSet(URI dataSchemaUri) throws IOException {

        // opens a stream to the URL passed as 
        // argument and reads its content into
        // an array of bytes.
        URL url = dataSchemaUri.toURL();
        try(InputStream stream = url.openStream()) {
            ByteArrayOutputStream output = new ByteArrayOutputStream();
            int nrBytesRead = 0;
            byte[] buffer = new byte[256];
            while((nrBytesRead = stream.read(buffer, 0, buffer.length)) != -1) {
                output.write(buffer, 0, nrBytesRead);
            }
            output.flush();
            
            buffer = output.toByteArray();

            return DescriptorProtos.FileDescriptorSet.parseFrom(buffer);
        }
    }

    /**
     * Converts the payload from its protobuf binary representation into a corresponding JSON
     * structure by leveraging the information supplied by the given file descriptor set.
     * 
     * @param descriptorSet     a {@link DescriptorProtos.FileDescriptorSet} instance that 
     *                          describes the schema of the entity peristed in the binary
     *                          protobuf.
     * 
     * @param messageType       a {@link String} representing the type of the root message
     *                          serialised in the bytes array.
     * 
     * @param protobuf          a {@literal byte} array containing the serialised version of
     *                          the protobuf entity.
     */
    protected JsonNode convertPayload(DescriptorProtos.FileDescriptorSet descriptorSet, String messageType, byte[] protobuf) throws IOException {

        return null;
    }


    /**
     * Generates a new event that matches the supplied event but whose content is replaced
     * with the new payload.
     * 
     * @param sourceEvent   a {@link CloudEvent} instance representing the source event to
     *                      copy all the attributes from. 
     * @param jsonPayload   a {@link JsonCloudEventData} instance representing the converted
     *                      payload of the original event into JSON format.
     * 
     * @return  a {@link CloudEvent} instance that contains the new payload expressed in JSON
     *          and all the original attributes of the supplied event.
     */
    protected CloudEvent repackageEvent(CloudEvent sourceEvent, JsonCloudEventData jsonPayload) {

        CloudEventBuilder builder = CloudEventBuilder.v1(sourceEvent)
                                        .withData("application/json", jsonPayload);

        return builder.build();
    }

}
