package org.cyverse.bench;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.FileInputStream;
import java.io.IOException;

import org.irods.jargon.core.connection.IRODSAccount;
import org.irods.jargon.core.connection.auth.AuthResponse;
import org.irods.jargon.core.exception.JargonException;
import org.irods.jargon.core.pub.IRODSFileSystem;
import org.irods.jargon.core.pub.io.IRODSFile;
import org.irods.jargon.core.pub.io.IRODSFileFactory;
import org.irods.jargon.core.pub.io.IRODSFileOutputStream;

/**
 * Hello world!
 *
 */
public class Upload 
{
    public static void main( String[] args )
    {
        String hostport = args[0];
        String user = args[1];
        String password = args[2];
        String zone = args[3];
        String localpath = args[4];
        String path = args[5];
        
        String host = hostport;
        int port = 1247;
        
        if (hostport.indexOf(":") >= 0) {
            String[] hostport_vars = hostport.split(":");
            host = hostport_vars[0].trim();
            port = Integer.parseInt(hostport_vars[1].trim());
        }

        if (port <= 0) {
            port = 1247;
        }

        try {
            IRODSAccount account = IRODSAccount.instance(host, port, user, password, "/" + zone, zone, "");
            IRODSFileSystem irodsFS = IRODSFileSystem.instance();

            AuthResponse authResponse = irodsFS.getIRODSAccessObjectFactory().authenticateIRODSAccount(account);

            if (!authResponse.isSuccessful()) {
                throw new IOException("cannot login");
            }

            int bufferSize = 8*1024*1024;


            IRODSFileFactory factory = irodsFS.getIRODSFileFactory(account);
            
            FileInputStream fis = new FileInputStream(localpath);
            BufferedInputStream bis = new BufferedInputStream(fis, bufferSize);
            
            IRODSFile ipath = factory.instanceIRODSFile(path);

            //IRODSFileOutputStream os = factory.instanceIRODSFileOutputStreamWithRerouting(ipath);
            IRODSFileOutputStream os = factory.instanceIRODSFileOutputStream(ipath);
            BufferedOutputStream bos = new BufferedOutputStream(os, bufferSize);

            byte[] buffer = new byte[bufferSize];

            while (true) {
                int readLen = bis.read(buffer);
                if (readLen > 0) {
                    bos.write(buffer, 0, readLen);
                } else {
                    break;
                }
            }

            bos.close();
            bis.close();

        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }
}
