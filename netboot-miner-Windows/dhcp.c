#include <stdio.h>
#include <winsock2.h>
#include <ws2tcpip.h>
#include <string.h>
#include <iphlpapi.h>

#pragma comment(lib, "iphlpapi.lib")
#pragma comment(lib, "Ws2_32.lib")

#define MALLOC(x) HeapAlloc(GetProcessHeap(), 0, (x))
#define FREE(x) HeapFree(GetProcessHeap(), 0, (x))


#define WORKING_BUFFER_SIZE 15000
#define MAX_TRIES 3

void getMyIp(u_long* ip, u_long* subnetmask)
{
    *ip = 0;
    *subnetmask = 0;
    // Declare and initialize variables.

    /* variables used for GetIfForwardTable */
    PMIB_IPFORWARDTABLE pIpForwardTable;
    DWORD dwSize = 0;
    DWORD dwRetVal = 0;

    char szDestIp[128];
    char szMaskIp[128];
    char szGatewayIp[128];

    struct in_addr IpAddr;

    int i;

    pIpForwardTable =
        (MIB_IPFORWARDTABLE *)MALLOC(sizeof(MIB_IPFORWARDTABLE));
    if (pIpForwardTable == NULL)
    {
        printf("Error allocating memory\n");
        return 1;
    }

    if (GetIpForwardTable(pIpForwardTable, &dwSize, 0) ==
        ERROR_INSUFFICIENT_BUFFER)
    {
        FREE(pIpForwardTable);
        pIpForwardTable = (MIB_IPFORWARDTABLE *)MALLOC(dwSize);
        if (pIpForwardTable == NULL)
        {
            printf("Error allocating memory\n");
            return 1;
        }
    }
    u_long mask = 0, nexthop = 0;
    /* Note that the IPv4 addresses returned in 
     * GetIpForwardTable entries are in network byte order 
     */
    if ((dwRetVal = GetIpForwardTable(pIpForwardTable, &dwSize, 0)) == NO_ERROR)
    {
        printf("\tNumber of entries: %d\n",
               (int)pIpForwardTable->dwNumEntries);
        for (i = 0; i < (int)pIpForwardTable->dwNumEntries; i++)
        {
            if ((u_long)pIpForwardTable->table[i].dwForwardDest != 0 && (u_long)pIpForwardTable->table[i].dwForwardMask != 0)
            {
                continue;
            }
            /* Convert IPv4 addresses to strings */
            IpAddr.S_un.S_addr =
                (u_long)pIpForwardTable->table[i].dwForwardDest;
            strcpy_s(szDestIp, sizeof(szDestIp), inet_ntoa(IpAddr));
            IpAddr.S_un.S_addr =
                (u_long)pIpForwardTable->table[i].dwForwardMask;
            mask = (u_long)pIpForwardTable->table[i].dwForwardMask;
            strcpy_s(szMaskIp, sizeof(szMaskIp), inet_ntoa(IpAddr));
            IpAddr.S_un.S_addr =
                (u_long)pIpForwardTable->table[i].dwForwardNextHop;
            nexthop = (u_long)pIpForwardTable->table[i].dwForwardNextHop;
            strcpy_s(szGatewayIp, sizeof(szGatewayIp), inet_ntoa(IpAddr));
            printf("mask,nexthop:%x,%x\n", mask, nexthop);
        }
        FREE(pIpForwardTable);
        // return 0;
    }
    else
    {
        printf("\tGetIpForwardTable failed.\n");
        FREE(pIpForwardTable);
        return 1;
    }
    /* Declare and initialize variables */

    // It is possible for an adapter to have multiple
    // IPv4 addresses, gateways, and secondary WINS servers
    // assigned to the adapter.
    //
    // Note that this sample code only prints out the
    // first entry for the IP address/mask, and gateway, and
    // the primary and secondary WINS server for each adapter.

    PIP_ADAPTER_INFO pAdapterInfo;
    PIP_ADAPTER_INFO pAdapter = NULL;

    /* variables used to print DHCP time info */
    struct tm newtime;
    char buffer[32];
    errno_t error;

    ULONG ulOutBufLen = sizeof(IP_ADAPTER_INFO);
    pAdapterInfo = (IP_ADAPTER_INFO *)MALLOC(sizeof(IP_ADAPTER_INFO));
    if (pAdapterInfo == NULL)
    {
        printf("Error allocating memory needed to call GetAdaptersinfo\n");
        return 1;
    }
    // Make an initial call to GetAdaptersInfo to get
    // the necessary size into the ulOutBufLen variable
    if (GetAdaptersInfo(pAdapterInfo, &ulOutBufLen) == ERROR_BUFFER_OVERFLOW)
    {
        FREE(pAdapterInfo);
        pAdapterInfo = (IP_ADAPTER_INFO *)MALLOC(ulOutBufLen);
        if (pAdapterInfo == NULL)
        {
            printf("Error allocating memory needed to call GetAdaptersinfo\n");
            return 1;
        }
    }

    if ((dwRetVal = GetAdaptersInfo(pAdapterInfo, &ulOutBufLen)) == NO_ERROR)
    {
        pAdapter = pAdapterInfo;
        while (pAdapter)
        {
            mask = inet_addr(pAdapter->IpAddressList.IpMask.String);
            u_long ipaddr = inet_addr(pAdapter->IpAddressList.IpAddress.String);
            if (mask != 0 && (nexthop & mask) == (ipaddr & mask))
            {

                printf("My Address:%s\n", pAdapter->IpAddressList.IpAddress.String);
                *ip = ipaddr;
                *subnetmask = mask;
                return;
            }
            pAdapter = pAdapter->Next;
        }
    }
    else
    {
        printf("GetAdaptersInfo failed with error: %d\n", dwRetVal);
    }
    if (pAdapterInfo)
        FREE(pAdapterInfo);

    return 0;
}

void DumpHex(const void* data, size_t size) {
	char ascii[17];
	size_t i, j;
	ascii[16] = '\0';
	for (i = 0; i < size; ++i) {
		printf("%02X ", ((unsigned char*)data)[i]);
		if (((unsigned char*)data)[i] >= ' ' && ((unsigned char*)data)[i] <= '~') {
			ascii[i % 16] = ((unsigned char*)data)[i];
		} else {
			ascii[i % 16] = '.';
		}
		if ((i+1) % 8 == 0 || i+1 == size) {
			printf(" ");
			if ((i+1) % 16 == 0) {
				printf("|  %s %d\n", ascii, i);
			} else if (i+1 == size) {
				ascii[(i+1) % 16] = '\0';
				if ((i+1) % 16 <= 8) {
					printf(" ");
				}
				for (j = (i+1) % 16; j < 16; ++j) {
					printf("   ");
				}
				printf("|  %s \n", ascii);
			}
		}
	}
}

int main()
{
    WSADATA wsaData;

    SOCKET sock;
    struct sockaddr_in addr;

    unsigned char buf[2048];

    WSAStartup(MAKEWORD(2, 0), &wsaData);

    sock = socket(AF_INET, SOCK_DGRAM, 0);
    // INT bAllow = 1;

    // setsockopt(sock, SOL_SOCKET, SO_BROADCAST,(char *)&bAllow, sizeof(bAllow));

    addr.sin_family = AF_INET;
    addr.sin_port = htons(67);
    addr.sin_addr.S_un.S_addr = INADDR_ANY;

    bind(sock, (struct sockaddr *)&addr, sizeof(addr));
    if(0 != WSAGetLastError()){
        printf("bind:%d\n", WSAGetLastError());
    }

    memset(buf, 0, sizeof(buf));
    struct sockaddr_in from;
    int from_size = sizeof(struct sockaddr_in);
    u_long myipaddr,mask;
    getMyIp(&myipaddr, &mask);

    for (;;)
    {

        int recvsize = recvfrom(sock, buf, sizeof(buf), 0, (struct sockaddr *)&from, &from_size);
        
        DumpHex(buf,recvsize);
        int flag = 0;

        for(int i=240; i < recvsize; ){
            if(buf[i] == 0xff){
                break;
            }
            unsigned char type = buf[i];
            i++;
            unsigned char size = buf[i];
            i++;
            char pxe[] = "PXEClient";
            char ipxe[] = "ipxe.pxe";

            if(type == 60 && strncmp(pxe,&buf[i],strlen(pxe)) == 0){
                printf("boot\n");//FIXME
                flag = 1;
                break;
            }

            i += size;
        }
        
        if(flag){
            int n = 0;
            buf[0] = 2;//BOOTREPLY
            n=20;

            (u_long) (*((u_long*)(&buf[n]))) = myipaddr;

            strcpy(&buf[108], "ipxe-default.pxe");
            memset(&buf[240],0,sizeof(buf)-240);
            n=240;
            buf[n] = 60;
            buf[n+1] = strlen("PXEClient");
            strcpy(&buf[n+2],"PXEClient");
            n += strlen("PXEClient")+2;
            buf[n] = 43;//encapsulated options
            buf[n+1] = 3;//size
            buf[n+2] = 6;//discovery control
            buf[n+3] = 1;//size
            buf[n+4] = 8;//flag
            n += 5;
            buf[n] = 0xff;//end
            n += 1;
            struct sockaddr_in bcast_address;
            memset(&bcast_address, 0, sizeof(bcast_address));
            bcast_address.sin_port = ntohs(68);
            bcast_address.sin_family = AF_INET;
            // inet_pton(AF_INET, "192.168.0.255", &bcast_address.sin_addr.s_addr);
            bcast_address.sin_addr.s_addr = (myipaddr & (mask)) | ~(mask);
            DumpHex(buf,n);
            sendto(sock, buf, n, 0, (struct sockaddr*)&bcast_address, sizeof(bcast_address));
        }
    }
    closesocket(sock);

    WSACleanup();

    return 0;
}

