package com.github.hyp0th3rmi4.protobuf.dynamic;


/**
 * JsonConverter converts the payload of the cloud event
 * that is encoded in Protobuf into a corresponding JSON
 * structure by relying upon the information stored in the
 * associated file descriptor.
 */
public class JsonConverter {

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

    }

}
