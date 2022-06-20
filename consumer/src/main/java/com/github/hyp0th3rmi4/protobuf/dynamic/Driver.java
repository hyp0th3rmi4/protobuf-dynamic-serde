package com.github.hyp0th3rmi4.protobuf.dynamic;

import org.apache.commons.cli.CommandLine;
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
            CommandLine cmdLine = parser.parser(options, commandLineArguments);

            // extract parameter values.
            String sourcePath = cmdLine.getOptionValue(Driver.SOURCE_PATH_OPTION);
            String targetPath = cmdLine.hasOption(Driver.TARGET_PATH_OPTION) ? cmdLine.getOptionValue(Driver.TARGET_PATH_OPTION) : null;


            JsonConverter converter = new JsonConverter();
            converter.convert(sourcePath, targetPath);

        } catch(ParseException pex) {

            System.out.println("Error:")
            System.out.println("Error while parsing arguments: " + pex.getMessage())
            System.out.println();

            Driver.printUsage(options);
            Driver.printHelp(options);

            exitCode = 1;
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
            .desc("Path to the file that will store the payload of the event converted from protobuf to JSON (if omitted the JSON is dumped to the console)."))
            .build();

        final Options options = new Options();
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
