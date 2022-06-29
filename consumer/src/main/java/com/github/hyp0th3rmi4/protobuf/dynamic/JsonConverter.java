package com.github.hyp0th3rmi4.protobuf.dynamic;

import java.io.IOException;
import java.io.InputStream;
import java.io.FileWriter;
import java.io.PrintWriter;
import java.io.ByteArrayOutputStream;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import static com.google.protobuf.DescriptorProtos.DescriptorProto;
import static com.google.protobuf.DescriptorProtos.FileDescriptorProto;
import static com.google.protobuf.DescriptorProtos.FileDescriptorSet;
import static com.google.protobuf.Descriptors.Descriptor;
import static com.google.protobuf.Descriptors.FieldDescriptor;
import com.google.protobuf.DynamicMessage;
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
     * A {@link ObjectMapper} instance used to both serialise the 
     * the {@link CloudEvent} instance into a JSON pretty printed
     * string, and create the JSON structure for the payload.
     */
    protected ObjectMapper mapper = new ObjectMapper();


    /**
     * An instance of {@link Writer} that is used to output messages
     * while processing the protobuf definition.
     */
    protected PrintWriter output = new PrintWriter(System.out);



    /**
     * Converts the binary representation of a protobuf message into
     * a corresponding JSON structure.
     * 
     * @param sourcePath    a {@link String} representing the path to
     *                      the file storing the serialised protobuf.
     * 
     * @param targetPath    a {@link String} representing the path to
     *                      the location where to store the serialised
     *                      JSON representation of the message. If set
     *                      to {@literal null} the JSON structure is
     *                      dumped to the standard output.
     */
    public void convertFromRaw(String sourcePath, String targetPath, String schemaUri) throws Exception {

        URI uri = new URI(schemaUri);
        FileDescriptorSet descriptorSet = this.resolveDescriptorSet(uri);

        byte[] protobuf = Files.readAllBytes(Paths.get(sourcePath));
        JsonNode message = this.convertMessage(descriptorSet, uri.getFragment(), protobuf);
        this.writeToTarget(message, targetPath);

    }

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
     * 
     * @param isRaw         a boolean flag indicating whether the source has
     *                      to be interpreted as a raw protobuf binary of the
     *                      serialised message (i.e. {@literal true}) or as
     *                      the JSON representation of a CloudEvent instance
     *                      whose base64 payload is the serialised protobuf
     *                      binary of the message.
     */
    public void convertFromCloudEvent(String sourcePath, String targetPath) throws Exception {

        // get the CloudEvebt and the descriptor set
        CloudEvent cloudEvent = this.getCloudEvent(sourcePath);

        System.out.println(cloudEvent);

        FileDescriptorSet descriptorSet = this.resolveDescriptorSet(cloudEvent.getDataSchema());

        // unpack the binary protobuf and convert it
        // to a corresponding payload expressed in 
        // JSON.
        String messageType = cloudEvent.getDataSchema().getFragment();
        byte[] protobuf = cloudEvent.getData().toBytes();
        JsonNode message = this.convertMessage(descriptorSet, messageType, protobuf);
        JsonCloudEventData newData = JsonCloudEventData.wrap(message);

        // regenerate the Cloud Event with the new payload.
        
        CloudEvent jsonCloudEvent = this.repackageEvent(cloudEvent, newData);


        this.writeToTarget(jsonCloudEvent, targetPath);
        

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
    protected FileDescriptorSet resolveDescriptorSet(URI dataSchemaUri) throws IOException {

        // opens a stream to the URL passed as 
        // argument and reads its content into
        // an array of bytes.
        URL url = dataSchemaUri.toURL();
        try(InputStream stream = url.openStream()) {
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            int nrBytesRead = 0;
            byte[] buffer = new byte[256];
            while((nrBytesRead = stream.read(buffer, 0, buffer.length)) != -1) {
                baos.write(buffer, 0, nrBytesRead);
            }
            baos.flush();
            
            buffer = baos.toByteArray();

            return FileDescriptorSet.parseFrom(buffer);
        }
    }

    /**
     * Serialises as a JSON document the given entity and saves it to the specified 
     * filePath if specified, otherwise it dumps the output to the standard output.
     * 
     * @param entity        a {@link Object} instance representing the entity  to be 
     *                      serialised as JSON. This can either be an instance of {@link 
     *                      CloudEvent} or a {@link JsonNode}.
     * 
     * @param targetPath    a {@link String} pointing to the file system location where 
     *                      to persist the JSON representation of <i>entity</i>. If not 
     *                      specified the output is the console.
     */
    protected void writeToTarget(Object entity, String targetPath) throws IOException {

        
        ObjectWriter jsonWriter = this.mapper.writerWithDefaultPrettyPrinter();
        if (targetPath == null) {
            try (FileWriter writer = new FileWriter(targetPath, true)) {

                jsonWriter.writeValue(writer, entity);
            }
        } else {    
            
            jsonWriter.writeValue(this.output, entity);
        }
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

    /**
     * Converts the message from its protobuf binary representation into a corresponding JSON
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
    protected JsonNode convertMessage(FileDescriptorSet descriptorSet, String messageType, byte[] protobuf) throws IOException {


        Descriptor type = this.findTypeForMessage(descriptorSet, messageType);
        if (type == null) {
            throw new IOException("Could not find message type: " + messageType);
        }
        
        DynamicMessage message = DynamicMessage.parseFrom(type, protobuf);
        ObjectNode root = this.mapper.createObjectNode();
        
        Map<FieldDescriptor, Object> fields = message.getAllFields();
        for(Map.Entry<FieldDescriptor,Object> entry : fields.entrySet()) {
            FieldDescriptor field = entry.getKey();
            Object value = this.getValueForField(field, entry.getValue());
            String name = field.getName();
            this.setNodeProperty(root, name, value);
        }

        return root;
    }  
    
   /**
    * Scans the given descriptor set to retrieve the {@link Descriptor} instance representing
    * the type specified by <i>messageType</i>. If there is no matching message type it does
    * return {@literal null}.
    * 
    * @param fileDescriptorSet     a {@link FileDescriptorSet} representing the set of file 
    *                              descriptor associated to the protobuf definitoin that is
    *                              representing the schema of the message.
    * @param messageType           a {@link String} representing the type (or name) of the
    *                              message that is serialised in the protobuf binary.
    * 
    * @return  a {@link Descriptor} instance describing the given <i>messageType</i> or 
    *          {@literal null} if not found.
    */
   protected Descriptor findTypeForMessage(FileDescriptorSet fileDescriptorSet, String messageType) {

       Descriptor descriptor = null;

       for (FileDescriptorProto fileDescriptor : fileDescriptorSet.getFileList()) {

           List<DescriptorProto> messages = fileDescriptor.getMessageTypeList();
           for(DescriptorProto message : messages) {
               
               if (message.getName().equals(messageType)) {
                   descriptor = message.getDescriptor();
               }
           }
       }
       return descriptor;
   }

  

    /**
     * Creates a {@link JsonNode} that represents the given value. This method
     * uses the supplied metadata to determine what type of node to create.
     * 
     * @param metadata  a {@link FieldDescriptor} containing the metadata for 
     *                  the node.
     * @param value     the value for the attribute represented by the given
     *                  metadata.
     */
    protected Object getValueForField(FieldDescriptor metadata, Object value) {

        Object fieldValue = null;

        if (metadata.isRepeated()) {

            ArrayNode array = this.mapper.createArrayNode();
            // the value is then an iterable
            Iterable<?> items = (Iterable<?>) value; 
            for(Object item : items) {
                Object itemValue = this.getNonRepeatedValue(metadata, item);
                if (item == null) {
                    array.addNull();
                } else {
                    switch(metadata.getJavaType()) {
                        case BOOLEAN:
                            array.add((boolean) itemValue);
                        break;
                        case BYTE_STRING:
                            array.add((byte[]) itemValue);
                        break;
                        case DOUBLE:
                            array.add((double) itemValue);
                        break;
                        case INT:
                            array.add((int) itemValue);
                        break;
                        case FLOAT:
                            array.add((float) itemValue);
                        break;
                        case LONG:
                            array.add((long) itemValue);
                        break;
                        case STRING:
                            array.add((String) itemValue);
                        break;
                        case MESSAGE:
                            array.add((JsonNode) itemValue);
                        break;
                        default:
                            this.output.println("<-- SKIP --> Unexpected type for array (type: " + itemValue.getClass().getCanonicalName() + ").");
                        break;
                    }
                }
            }
            fieldValue = array;
        } else {

            fieldValue = this.getNonRepeatedValue(metadata, value);

        }
        return fieldValue;
    }

    /**
     * Returns the value mapped by the given <i>metadata</i> and <i>value</i>. This
     * implementation does validate the type of the value based on the information
     * stored in the given metadata descriptor and produces the corresponding object
     * to be added in the JSON tree that is being built to convert the original 
     * protobuf representation. Simple types are (int, string, bool, double, bytes,
     * and other primitive types) are simply passed through, while maps, messages
     * and groups are unwrapped.
     * 
     * @param metadata      a {@link FieldDescriptor} implementation that is used to
     *                      carry metadata about the field definition that matches
     *                      the supplied value.
     * 
     * @param value         an {@link Object} instance that represents the value of
     *                      the attribute represented by the descriptor.
     * 
     * @return an {@link Object} intance representing the value mapped by the method.
     */
    protected Object getNonRepeatedValue(FieldDescriptor metadata, Object value) {

        Object v = null;
        if (metadata.isMapField()) {
            ObjectNode map = this.mapper.createObjectNode();
            this.output.println("<MAP> " + value);
            v = map;
        } else {
            switch(metadata.getType()) {
                case ENUM:
                    this.output.println("<ENUM>: " + value);
                break;
                case MESSAGE:
                    this.output.println("<MESSAGE>: " + value);
                break;
                case GROUP:
                    this.output.println("<GROUP>: " + value);
                break;
                default:
                    v = value;
                break;
            }
        } 

        return v;
    }

    /**
     * Invokes the appropriate setter according to the concrete type of the given
     * value. This method is needed because {@link ObjectNode} does not have a
     * generic setter method for arguments of type {@link Object}.
     * 
     * @param node      an instance of {@link ObjectNode} representing the container
     *                  of the property.
     * @param name      a {@link String} representing the name of the property.
     * @param value     a {@link Object} representing the value of the property.
     */
    protected void setNodeProperty(ObjectNode node, String name, Object value) {

        if (value != null) {
            if (value instanceof JsonNode) {
                node.put(name, (JsonNode) value);
            } else {
                switch (value.getClass().getCanonicalName()) {
                    case "java.lang.Boolean":
                        node.put(name, (boolean) value);
                    break;
                    case "byte[]":
                        node.put(name, (byte[]) value);
                    break;
                    case "java.lang.Double":
                        node.put(name, (double) value);
                    break;
                    case "java.lang.Integer":
                        node.put(name, (int) value);
                    break;
                    case "java.lang.Float":
                        node.put(name, (float) value);
                    break;
                    case "java.lang.Long":
                        node.put(name, (long) value);
                    break;
                    case "java.lang.String":
                        node.put(name, (String) value);
                    break;
                    default:
                        this.output.println("<-- SKIP --> " + name + " (type: " + value.getClass().getCanonicalName() + ")");
                }
            }
        } else {
            node.putNull(name);
        }
    }
}
