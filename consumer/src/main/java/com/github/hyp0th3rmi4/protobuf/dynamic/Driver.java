package com.github.hyp0th3rmi4.protobuf.dynamic;


import java.io.PrintWriter;
import java.text.ParseException;

import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.CommandLineParser;
import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.Options;



/**
 * Driver is the main class of the package and implements the
 * necessary tasks to parse the command line arguments and then
 * trigger an instance of the {@link JsonConverter} to convert
 * the content of the cloud event from protobuf to JSON.
 */
public class Driver {

    /**
     * Long name of the command line option that indicates the
     * source format of the content to interpret and deserialise.
     */
    private static final String RAW_OPTION  = "raw";

    /**
     * Long name of the command line option that stores a reference
     * to the URI of the schema used to interpret the protobuf
     * binary representing the message.
     */
    private static final String SCHEMA_URI_OPTION = "schema_uri";

    /**
     * Long name of the command line option storing the path to the
     * file that contains the cloudevent instance encoded in JSON
     * format.
     */
    private static final String SOURCE_PATH_OPTION = "source_path";
    /**
     * Long name of the command line option storing the path to the
     * file that will be created/overwritten to store the payload
     * of the event stored in JSON format.
     */
    private static final String TARGET_PATH_OPTION = "target_path";

    /**
     * A {@link String} representing the identifier of the syntax used to
     * show the help and usage text.
     */
    private static final String SYNTAX = "Main";

    /**
     * A {@link HelpFormatter} instance used to print usage and help text
     * to the standard output in case of need.
     */
    private static HelpFormatter formatter = new HelpFormatter();
    

    /**
     * Main driver of the application. Reads and parses the command line arguments
     * and executes the conversion of the payload of the event. If there is any
     * error it dumps the error on the standard output.
     */
    public static void main(String[] args) {

        int exitCode = 0;

        final Options options = Driver.newOptions();
        final CommandLineParser parser = new DefaultParser();
        try {

            // parse command line arguments
            CommandLine cmdLine = parser.parse(options, args);

            // extract parameter values.
            String sourcePath = cmdLine.getOptionValue(Driver.SOURCE_PATH_OPTION);
            String targetPath = cmdLine.hasOption(Driver.TARGET_PATH_OPTION) ? cmdLine.getOptionValue(Driver.TARGET_PATH_OPTION) : null;
            boolean isRaw = cmdLine.hasOption(Driver.RAW_OPTION);

            JsonConverter converter = new JsonConverter();
            if (isRaw) {

                if (cmdLine.hasOption(Driver.SCHEMA_URI_OPTION)) {

                    String schemaUri = cmdLine.getOptionValue(Driver.SCHEMA_URI_OPTION);
                    converter.convertFromRaw(sourcePath, targetPath, schemaUri);
                
                } else {
                    throw new Exception("Invalid arguments: --raw (-r) was specified without a schema URI.");
                }

            } else  {

                converter.convertFromCloudEvent(sourcePath, targetPath);
            }

        } catch(ParseException pex) {

            System.out.println("Error:");
            System.out.println("Error while parsing arguments: " + pex.getMessage());
            System.out.println();

            Driver.showUsage(options);
            Driver.showHelp(options);

            exitCode = 1;
        
        } catch(Exception ex) {

            System.out.println("Error:");
            System.out.println("Type: " + ex.getClass().getName());
            System.out.println("Message: " + ex.getMessage());
            System.out.println();

            exitCode = 2;
        }

        System.exit(exitCode);
    }

    /**
     * Creates command line parsing options used to validate and
     * extract the information from the command line arguments
     * passed to the process.
     * 
     * @return  an instance of {@link Options} representing the 
     *          expected configuration optins for the command line.
     */
    private static Options newOptions() {

        final Option schemaUri = Option.builder("u")
            .required(false)
            .hasArg(true)
            .longOpt(Driver.SCHEMA_URI_OPTION)
            .desc("URI to the schema definition for the message stored in the protobuf binary. This parameter is only considered when the --raw option is passed, otherwise the 'dataschema' attribute will be considered when the source format is a cloud event.")
            .build();

        final Option raw = Option.builder("r")
            .required(false)
            .hasArg(false)
            .longOpt(Driver.RAW_OPTION)
            .desc("If specified the source file is interpreted as a raw protobuf binary representing the message, rather than a JSON representation of a CloudEvent with a base64 protobuf binary for payload.")
            .build();

        final Option sourcePath = Option.builder("s")
            .required(true)
            .hasArg(true)
            .longOpt(Driver.SOURCE_PATH_OPTION)
            .desc("Path to the file that contains the cloud event in JSON format.")
            .build();

        final Option targetPath = Option.builder("t")
            .required(false)
            .hasArg(true)
            .longOpt(Driver.TARGET_PATH_OPTION)
            .desc("Path to the file that will store the payload of the event converted from protobuf to JSON (if omitted the JSON is dumped to the console).")
            .build();

        final Options options = new Options();
        options.addOption(raw);
        options.addOption(schemaUri);
        options.addOption(sourcePath);
        options.addOption(targetPath);

        return options;
    }

    /**
     * Prints to the standard output the usage information
     * associated to the application.
     * 
     * @param options   a {@link Options} implementation 
     *                  defining the parameters that are
     *                  expected/accepted by the application.
     */
    private static void showHelp(final Options options) {

        System.out.println("Help");
        Driver.formatter.printHelp(Driver.SYNTAX, null, options, null);
        System.out.println();
    }

    /**
     * Prints to the standard output the help information
     * associated to the application.
     * 
     * @param options   a {@link Options} implementation 
     *                  defining the parameters that are
     *                  expected/accepted by the application.
     */
    private static void showUsage(final Options options) {

        System.out.println("Usage:");
        final PrintWriter pw = new PrintWriter(System.out);
        Driver.formatter.printUsage(pw, 80, Driver.SYNTAX, options);
        pw.flush();
        System.out.println();
    }
}
