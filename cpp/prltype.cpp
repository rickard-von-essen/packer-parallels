#include "ParallelsVirtualizationSDK/Parallels.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

PRL_RESULT LoginLocal(PRL_HANDLE &hServer) {

    PRL_HANDLE hJob, hJobResult = PRL_INVALID_HANDLE;
    PRL_RESULT err, nJobReturnCode = PRL_ERR_UNINITIALIZED;

    err = PrlApi_InitEx(PARALLELS_API_VER, PAM_DESKTOP, 0, 0);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlApi_InitEx returned with error: %s.\n", prl_result_to_string(err));
        PrlApi_Deinit();
        return err;
    }
    err = PrlSrv_Create(&hServer);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlSvr_Create failed, error: %s", prl_result_to_string(err));
        PrlApi_Deinit();
        return err;
    }
    hJob = PrlSrv_LoginLocal(hServer, NULL, NULL, PSL_NORMAL_SECURITY);
    err = PrlJob_Wait(hJob, 1000);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlJob_Wait for PrlSrv_Login returned with error: %s\n", prl_result_to_string(err));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
    }
    err = PrlJob_GetRetCode(hJob, &nJobReturnCode);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlJob_GetRetCode returned with error: %s\n", prl_result_to_string(err));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
    }
    if (PRL_FAILED(nJobReturnCode)) {
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        fprintf(stderr, "Login job returned with error: %s\n", prl_result_to_string(nJobReturnCode));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
    }
    return PRL_ERR_SUCCESS;
}

PRL_RESULT LogOff(PRL_HANDLE &hServer) {
    PRL_HANDLE hJob, hJobResult = PRL_INVALID_HANDLE;
    PRL_RESULT err, nJobReturnCode = PRL_ERR_UNINITIALIZED;

    nanosleep((struct timespec[]){{1, 500000000}}, NULL);

    hJob = PrlSrv_Logoff(hServer);
    err = PrlJob_Wait(hJob, 1000);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlJob_Wait for PrlSrv_Logoff returned error: %s\n", prl_result_to_string(err));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
     }
    err = PrlJob_GetRetCode(hJob, &nJobReturnCode);
    if (PRL_FAILED(err)) {
        fprintf(stderr, "PrlJob_GetRetCode failed for PrlSrv_Logoff with error: %s\n", prl_result_to_string(err));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
    }
    if (PRL_FAILED(nJobReturnCode)) {
        fprintf(stderr, "PrlSrv_Logoff failed with error: %s\n", prl_result_to_string(nJobReturnCode));
        PrlHandle_Free(hJob);
        PrlHandle_Free(hServer);
        PrlApi_Deinit();
        return err;
    }
    PrlHandle_Free(hJob);
    PrlHandle_Free(hServer);
    PrlApi_Deinit();
    return PRL_ERR_SUCCESS;
}

PRL_RESULT GetVm(PRL_HANDLE &hRetVm, PRL_HANDLE &hServer, char* vmName ) {
    PRL_HANDLE hJob, hJobResult  = PRL_INVALID_HANDLE;
    PRL_RESULT ret, nJobReturnCode  = PRL_ERR_UNINITIALIZED;
    hJob = PrlSrv_GetVmList(hServer);
    ret = PrlJob_Wait(hJob, 10000);
    if (PRL_FAILED(ret)) {
        fprintf(stderr, "PrlJob_Wait for PrlSrv_GetVmList returned with error: %s\n",
        prl_result_to_string(ret));
        PrlHandle_Free(hJob);
        return ret;
    }
    ret = PrlJob_GetRetCode(hJob, &nJobReturnCode);
    if (PRL_FAILED(ret)) {
        fprintf(stderr, "PrlJob_GetRetCode returned with error: %s\n", prl_result_to_string(ret));
        PrlHandle_Free(hJob);
        return ret;
    }
    if (PRL_FAILED(nJobReturnCode)) {
        fprintf(stderr, "PrlSrv_GetVmList returned with error: %s\n",
        prl_result_to_string(ret));
        PrlHandle_Free(hJob);
        return ret;
    }
    ret = PrlJob_GetResult(hJob, &hJobResult);
    if (PRL_FAILED(ret)) {
        fprintf(stderr, "PrlJob_GetResult returned with error: %s\n", prl_result_to_string(ret));
        PrlHandle_Free(hJob);
        return ret;
    }
        PrlHandle_Free(hJob);
        PRL_UINT32 nParamsCount = 0;
        ret = PrlResult_GetParamsCount(hJobResult, &nParamsCount);
        for (PRL_UINT32 i = 0; i < nParamsCount; ++i)
        {
             PRL_HANDLE hVm = PRL_INVALID_HANDLE;
             PrlResult_GetParamByIndex(hJobResult, i, &hVm);
             char szVmNameReturned[1024];
             PRL_UINT32 nBufSize = sizeof(szVmNameReturned);
             ret = PrlVmCfg_GetName(hVm, szVmNameReturned, &nBufSize);
             if (PRL_FAILED(ret)) {
                fprintf(stderr, "PrlVmCfg_GetName returned with error (%s)\n",
                prl_result_to_string(ret));
             } else {
                if (strcmp(vmName, szVmNameReturned) == 0) {
                  hRetVm = hVm;
                  return PRL_ERR_SUCCESS;
                }
             }
             PrlHandle_Free(hVm);
         }
        fprintf(stderr, "No VM: \"%s\" was found!\n", vmName);
    return PRL_ERR_OPERATION_FAILED;
}

PRL_RESULT ConnectToConsole( PRL_HANDLE vm ) {

  PRL_RESULT err, nJobReturnCode = PRL_ERR_UNINITIALIZED;
  PRL_HANDLE hJob = PrlDevDisplay_ConnectToVm( vm, PDCT_LOW_QUALITY_WITHOUT_COMPRESSION );
  err = PrlJob_Wait(hJob, 10000);
  if (PRL_FAILED(err)) {
      fprintf(stderr, "PrlJob_Wait for PrlDevDisplay_ConnectToVm returned with error: %s\n",
      prl_result_to_string(err));
      PrlHandle_Free(hJob);
      return err;
  }
  err = PrlJob_GetRetCode(hJob, &nJobReturnCode);
  if (PRL_FAILED(err)) {
      fprintf(stderr, "PrlJob_GetRetCode returned with error: %s\n", prl_result_to_string(err));
      PrlHandle_Free(hJob);
      return err;
  }
  if (PRL_FAILED(nJobReturnCode)) {
      fprintf(stderr, "PrlDevDisplay_ConnectToVm returned with error: %s\n", prl_result_to_string(nJobReturnCode));
      PrlHandle_Free(hJob);
      return nJobReturnCode;
  }

  return PRL_ERR_SUCCESS;
}

PRL_RESULT DisconnectFromConsole( PRL_HANDLE vm ) {

  PRL_RESULT err = PrlDevDisplay_DisconnectFromVm( vm );
  if (PRL_FAILED(err)) {
      fprintf(stderr, "PrlDevDisplay_DisconnectFromVm returned with error: %s\n", prl_result_to_string(err));
  }
  return err;
}

PRL_RESULT SendScanCode( PRL_HANDLE vm, int scanCode ) {
  fprintf(stdout, " %d", scanCode);
  nanosleep((struct timespec[]){{0, 20000000}}, NULL);
  PRL_RESULT err;
  if (scanCode < 0x80)
    err = PrlDevKeyboard_SendKeyEvent( vm, scanCode, PKE_PRESS );
  else
    err = PrlDevKeyboard_SendKeyEvent( vm, scanCode - 0x80, PKE_RELEASE );
  if (PRL_FAILED(err)) {
      fprintf(stdout, "PrlDevKeyboard_SendKeyEvent returned with error: %s.\n", prl_result_to_string(err));
      PrlApi_Deinit();
      PrlHandle_Free( vm );
  }
  return err;
}

int main( int argc, char **argv ) {

  if (argc < 2) {
    fprintf(stdout, "Usage: prltype <vm-name> <scan-codes...>\nexample: prltype my-vm 30 158\n  This will press and release \"a\" on vm \"my-vm\".\n");
    return -1;
  }

  char vmName[100];
  strcpy(vmName, argv[1]);
  int scanCodes [argc - 1];
  for (int i = 2; i < argc; i++) {
    scanCodes[i - 2] = strtol(argv[i], NULL, 16);
    //scanCodes[i - 2] = atoi(argv[i]);
  }

  PRL_HANDLE hServer, hVm = PRL_INVALID_HANDLE;
  PRL_RESULT err;

  err = LoginLocal(hServer);
  if (PRL_FAILED(err))
    return err;
  err = GetVm(hVm, hServer, vmName);
  if (PRL_FAILED(err))
    return err;
  err = ConnectToConsole(hVm);
  if (PRL_FAILED(err))
    return err;

  fprintf(stdout, "Connected to console for VM: %s. Sending scancode(s) [", vmName);
  for (int i = 0; i < argc - 2; i++) {
    err = SendScanCode(hVm, scanCodes[i]);
    if (PRL_FAILED(err))
      return err;
  }
  fprintf(stdout, " ]\n\n");

  err = DisconnectFromConsole(hVm);
  if (PRL_FAILED(err))
    return err;

  LogOff( hServer );
}
